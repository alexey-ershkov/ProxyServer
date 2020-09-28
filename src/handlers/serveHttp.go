package handlers

import (
	"github.com/sirupsen/logrus"
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

func ServeHttp(respWriter http.ResponseWriter, request *http.Request) {

	logrus.Info("Request: " + request.RequestURI)

	var err error
	var hh Handler

	if request.Method == http.MethodConnect {
		hh, err = NewHttpsHandler(respWriter, request)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		hh = NewHttpHandler(respWriter, request)
	}

	err = hh.ProxyRequest()
	if err != nil {
		logrus.Error(err)
	}

	hh.Defer()
}
