package handlers

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
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
		proxyResp, err = handleHTTPS(request)
	} else {
		proxyResp, err = handleHTTP(request)
	}
	if err != nil {
		logrus.Error(err)
	}

	respWriter.WriteHeader(proxyResp.StatusCode)
	copyHeaders(proxyResp.Header, respWriter.Header())

	_, err = io.Copy(respWriter, proxyResp.Body)
	if err != nil {
		logrus.Error(err)
	}
}
