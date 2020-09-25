package handlers


import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func MainHandler(writer http.ResponseWriter, reader *http.Request){
	logrus.Debug("Handle",reader.Header, reader.Host, reader.GetBody)
	writer.WriteHeader(200)
	err := reader.Write(writer)

	if err != nil {
		logrus.Error(err)
	}
}
