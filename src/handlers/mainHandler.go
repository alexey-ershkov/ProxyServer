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
	for header, values := range from{
		for _, value := range values{
			to.Add(header, value)
		}
	}
}

func MainHandler(respWriter http.ResponseWriter, request *http.Request) {
	logrus.WithFields(logrus.Fields{
		"URL": request.RequestURI,
	}).Debug("Main Handler")

	proxyResp := handleHTTP(request)

	defer proxyResp.Body.Close()

	respWriter.WriteHeader(proxyResp.StatusCode)
	copyHeaders(proxyResp.Header, respWriter.Header())

	_, err := io.Copy(respWriter, proxyResp.Body)
	if err != nil {
		logrus.Error(err)
	}

}
