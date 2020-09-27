package handlers

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
)

func handleHTTPS(respWriter http.ResponseWriter, connect *http.Request) (*http.Response, error) {

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error(err)
	}

	genScriptDir := pwd + "/certGen"
	certsDir := pwd + "/certs/"

	parsedUrl, err := url.Parse(connect.RequestURI)
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

	config := new(tls.Config)

	config.Certificates = []tls.Certificate{cert}
	config.InsecureSkipVerify = true
	config.ServerName = parsedUrl.Scheme

	var serverConnection *tls.Conn
	var clientConnection net.Conn

	clientConnection, err = handshake(respWriter, config)
	if err != nil {
		logrus.Println("handshake", connect.Host, err)
		return nil, err
	}
	defer clientConnection.Close()

	serverConnection, err = tls.Dial("tcp", connect.Host, config)
	if err != nil {
		logrus.Fatal(err)
	}
	defer serverConnection.Close()

	reader := bufio.NewReader(clientConnection)
	request, err := http.ReadRequest(reader)
	if err != nil {
		logrus.Error(err)
	}

	fmt.Println(request)

	rawReq, err := httputil.DumpRequest(request, true)
	_, err = serverConnection.Write(rawReq)
	if err != nil {
		logrus.Error(err)
	}

	writer := bufio.NewReader(serverConnection)
	response, err := http.ReadResponse(writer, request)

	rawResp, err := httputil.DumpResponse(response, true)
	_, err = clientConnection.Write(rawResp)
	if err != nil {
		logrus.Error(err)
	}

	clientConnection.Close()
	serverConnection.Close()

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
