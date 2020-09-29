package proxyHandlers

import (
	"Proxy/db"
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

func ServeHttp(respWriter http.ResponseWriter, request *http.Request) {

	logrus.Info("Request: " + request.RequestURI)

	var err error
	var hh Handler

	dbConn, err := db.CreateNewDatabaseConnection()
	if err != nil {
		logrus.Warn("Can't connect to database")
		logrus.Fatal(err)
	}

	if request.Method == http.MethodConnect {
		hh, err = NewHttpsHandler(respWriter, request, dbConn)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		hh = NewHttpHandler(respWriter, request, dbConn)
	}

	err = hh.ProxyRequest()
	if err != nil {
		logrus.Error(err)
	}

	defer hh.Defer()
}
