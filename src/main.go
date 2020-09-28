package main

import (
	"Proxy/handlers"
	"github.com/sirupsen/logrus"
	"net/http"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

func main() {
	PORT := "1080"
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: http.HandlerFunc(handlers.ServeHttp),
	}

	logrus.Info("Server started on port ", PORT)

	logrus.Fatal(server.ListenAndServe())
}
