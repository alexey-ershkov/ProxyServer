package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

func contains(slice []string, value string) bool {
	for _, i := range slice {
		if value == i {
			return true
		}
	}
	return false
}

func copyHeaders(from, to http.Header) {
	for header, values := range from {
		for _, value := range values {
			to.Add(header, value)
		}
	}
}

func MainHandler(respWriter http.ResponseWriter, request *http.Request) {

	logrus.Info("Request: " + request.RequestURI)

	var proxyResp *http.Response
	var err error

	if request.Method == http.MethodConnect {
		parsedUrl, err := url.Parse(request.RequestURI)
		if err != nil {
			logrus.Error(err)
		}

		config, err := getHttpsConfig(parsedUrl.Scheme)

		clientTcpSocket, err := setupHttpsClientConnection(respWriter, config)
		if err != nil {
			logrus.Error(err)
		}

		serverTcpSocket, err := setupHttpsServerConnection(request.Host, config)

		request, err := getHttpsRequest(clientTcpSocket)

		proxyResp, err = doHttpsProxyRequest(request, serverTcpSocket)
		err = sendHttpsProxyResponse(proxyResp, clientTcpSocket)

		serverTcpSocket.Close()
		clientTcpSocket.Close()

	} else {
		proxyResp, err = doHttpProxyRequest(request)
		sendHttpProxyResponse(respWriter, proxyResp)
	}
	if err != nil {
		logrus.Error(err)
	}

}
