package main

import (
	"Proxy/handlers"
	"github.com/sirupsen/logrus"
	"net/http"
)



func main()  {
	logrus.SetLevel(logrus.InfoLevel)

	server := &http.Server{
		Addr: ":3000",
		Handler: http.HandlerFunc(handlers.MainHandler),
	}

	logrus.Info("Server started on port 3000")

	logrus.Fatal(server.ListenAndServe())
}
