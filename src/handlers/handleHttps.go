package handlers

import (
	"bufio"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
)

func doHttpsProxyRequest(request *http.Request, serverConnection *tls.Conn) (*http.Response, error) {



	rawReq, err := httputil.DumpRequest(request, true)
	_, err = serverConnection.Write(rawReq)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewReader(serverConnection)
	response, err := http.ReadResponse(writer, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func setupHttpsClientConnection(respWriter http.ResponseWriter, config *tls.Config) (net.Conn, error) {
	raw, _, err := respWriter.(http.Hijacker).Hijack()
	if err != nil {
		return nil, err
	}

	_, err = raw.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		raw.Close()
		return nil, err
	}

	clientConnection := tls.Server(raw, config)
	err = clientConnection.Handshake()
	if err != nil {
		clientConnection.Close()
		raw.Close()
		return nil, err
	}

	return clientConnection, nil
}

func setupHttpsServerConnection(host string, config *tls.Config) (*tls.Conn, error) {
	serverConnection, err := tls.Dial("tcp", host, config)
	if err != nil {
		return nil, err
	}

	return serverConnection, nil
}

func getHttpsRequest(clientConnection net.Conn) (*http.Request, error) {
	reader := bufio.NewReader(clientConnection)
	request, err := http.ReadRequest(reader)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func sendHttpsProxyResponse(response *http.Response, clientConnection net.Conn) error {
	rawResp, err := httputil.DumpResponse(response, true)
	_, err = clientConnection.Write(rawResp)
	if err != nil {
		return err
	}

	return nil
}
