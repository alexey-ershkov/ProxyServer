package proxyHandlers

import (
	"Proxy/src/db"
	"Proxy/src/models"
	"bufio"
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type HttpsHandler struct {
	connectRequest      *http.Request
	clientRequest       *http.Request
	respWriter          http.ResponseWriter
	config              *tls.Config
	serverTcpConnection *tls.Conn
	clientTcpConnection net.Conn
	parsedUrl           *url.URL
	proxyResp           *http.Response
	dbConn              *db.Database
}

func NewHttpsHandler(respWriter http.ResponseWriter, connectRequest *http.Request, dbConn *db.Database) (*HttpsHandler, error) {
	hh := &HttpsHandler{}
	hh.respWriter = respWriter
	hh.connectRequest = connectRequest
	hh.dbConn = dbConn

	var err error

	hh.parsedUrl, err = url.Parse(connectRequest.RequestURI)
	if err != nil {
		return nil, err
	}

	err = hh.setupHttpsConfig()
	if err != nil {
		return nil, err
	}

	err = hh.setupHttpsClientConnection()
	if err != nil {
		return nil, err
	}

	err = hh.setupHttpsServerConnection()
	if err != nil {
		return nil, err
	}

	return hh, nil
}

func (hh *HttpsHandler) ProxyRequest() error {

	var err error

	err = hh.getHttpsRequest()
	if err != nil {
		return err
	}

	reqDump, err := httputil.DumpRequest(hh.clientRequest, true)
	if err != nil {
		logrus.Warn("Can't dump request")
	}

	dbReq := models.DatabaseReq{
		Host:    "https://" + hh.parsedUrl.Scheme + hh.clientRequest.URL.Path,
		Request: string(reqDump),
	}

	hh.dbConn.InsertRequest(dbReq)

	err = hh.doHttpsProxyRequest()
	if err != nil {
		return err
	}

	err = hh.sendHttpsProxyResponse()
	if err != nil {
		return err
	}

	return nil
}

func (hh *HttpsHandler) Defer() {
	hh.clientTcpConnection.Close()
	hh.serverTcpConnection.Close()
	hh.dbConn.Close()
}

func (hh *HttpsHandler) doHttpsProxyRequest() error {

	rawReq, err := httputil.DumpRequest(hh.clientRequest, true)
	_, err = hh.serverTcpConnection.Write(rawReq)
	if err != nil {
		return err
	}

	writer := bufio.NewReader(hh.serverTcpConnection)
	response, err := http.ReadResponse(writer, hh.clientRequest)
	if err != nil {
		return err
	}

	hh.proxyResp = response

	return nil
}

func (hh *HttpsHandler) setupHttpsClientConnection() error {
	raw, _, err := hh.respWriter.(http.Hijacker).Hijack()
	if err != nil {
		return err
	}

	_, err = raw.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		raw.Close()
		return err
	}

	clientConn := tls.Server(raw, hh.config)
	err = clientConn.Handshake()
	if err != nil {
		clientConn.Close()
		raw.Close()
		return err
	}

	hh.clientTcpConnection = clientConn

	return nil
}

func (hh *HttpsHandler) setupHttpsServerConnection() error {
	serverConnection, err := tls.Dial("tcp", hh.connectRequest.Host, hh.config)
	if err != nil {
		return err
	}

	hh.serverTcpConnection = serverConnection

	return nil
}

func (hh *HttpsHandler) getHttpsRequest() error {
	reader := bufio.NewReader(hh.clientTcpConnection)
	request, err := http.ReadRequest(reader)
	if err != nil {
		return err
	}

	hh.clientRequest = request

	return nil
}

func (hh *HttpsHandler) sendHttpsProxyResponse() error {
	rawResp, err := httputil.DumpResponse(hh.proxyResp, true)
	_, err = hh.clientTcpConnection.Write(rawResp)
	if err != nil {
		return err
	}

	return nil
}

func (hh *HttpsHandler) setupHttpsConfig() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	genScriptAndRootCaDir := pwd + "/certGen"
	certsDir := pwd + "/certs/"

	certFilename := certsDir + hh.parsedUrl.Scheme + ".crt"

	_, err = os.Stat(certFilename)
	if os.IsNotExist(err) {
		err = genProxyCert(genScriptAndRootCaDir, "/gen_cert.sh", hh.parsedUrl.Scheme, certsDir)
		if err != nil {
			return err
		}
	}

	cert, err := tls.LoadX509KeyPair(certFilename, genScriptAndRootCaDir+"/cert.key")
	if err != nil {
		return err
	}

	config := new(tls.Config)
	config.Certificates = []tls.Certificate{cert}
	config.ServerName = hh.parsedUrl.Scheme

	hh.config = config

	return nil
}
