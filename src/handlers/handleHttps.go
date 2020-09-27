package handlers

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

func handleHTTPS(respWriter http.ResponseWriter, request *http.Request) (*http.Response, error) {

	var serverConnection *tls.Conn
	var clientConnection net.Conn

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error(err)
	}

	genScriptDir := pwd + "/certGen"
	certsDir := pwd + "/certs/"

	parsedUrl, err := url.Parse(request.RequestURI)
	if err != nil {
		logrus.Error(err)
	}

	genCmd := exec.Command(genScriptDir+"/gen_cert.sh", parsedUrl.Scheme, "01", genScriptDir, certsDir)
	info, err := genCmd.CombinedOutput()
	if err != nil {
		logrus.Error(string(info), err)
	} else {
		logrus.Debug(string(info))
	}

	cert, err := tls.LoadX509KeyPair(certsDir+parsedUrl.Scheme+".crt", genScriptDir+"/cert.key")
	if err != nil {
		logrus.Error(err)
	}

	//certCa, err := tls.LoadX509KeyPair(genScriptDir+"/ca.crt", genScriptDir+"/ca.key")
	//if err != nil {
	//	logrus.Error(err)
	//}

	serverConfig := new(tls.Config)
	//clientConfig := new(tls.Config)

	certPool := x509.NewCertPool()
	readCert,err := ioutil.ReadFile(certsDir+parsedUrl.Scheme+".crt")
	if err != nil {
		logrus.Fatal(err)
	}
	readCertCA,err := ioutil.ReadFile(genScriptDir+"/ca.crt")
	if err != nil {
		logrus.Fatal(err)
	}
	certPool.AppendCertsFromPEM(readCert)
	certPoolCA := x509.NewCertPool()
	certPoolCA.AppendCertsFromPEM(readCertCA)


	serverConfig.ClientCAs = certPool
	serverConfig.RootCAs = certPoolCA
	serverConfig.Certificates = []tls.Certificate{cert}
	serverConfig.InsecureSkipVerify = true
	serverConfig.ServerName = parsedUrl.Scheme


	serverConnection, err = tls.Dial("tcp", request.Host, serverConfig)
	if err != nil {
		logrus.Fatal(err)
	}


	clientConnection, err = handshake(respWriter, serverConfig)
	if err != nil {
		logrus.Println("handshake", request.Host, err)
		return nil, err
	}

	if clientConnection != nil && serverConnection != nil {
		logrus.Fatal("Ok")
	}

	return nil, nil
}

var okHeader = []byte("HTTP/1.1 200 Connection established\r\n\r\n")

func handshake(w http.ResponseWriter, config *tls.Config) (net.Conn, error) {
	raw, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, "no upstream", 503)
		return nil, err
	}
	if _, err = raw.Write(okHeader); err != nil {
		raw.Close()
		return nil, err
	}
	conn := tls.Server(raw, config)
	err = conn.Handshake()
	if err != nil {
		conn.Close()
		raw.Close()
		return nil, err
	}
	return conn, nil
}
