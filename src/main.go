package main

import (
	"Proxy/proxyHandlers"
	"Proxy/repeatReqHandlers"
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

	ProxyPort := "1080"
	RepeatPort := "80"

	server := &http.Server{
		Addr:    ":" + ProxyPort,
		Handler: http.HandlerFunc(proxyHandlers.ServeHttp),
	}

	http.HandleFunc("/", repeatReqHandlers.SendRequestList)
	http.HandleFunc("/req", repeatReqHandlers.ExecRepReq)

	logrus.Info("Proxy server started on port ", ProxyPort)
	logrus.Info("Repeat server started on port ", RepeatPort)

	go http.ListenAndServe(":"+RepeatPort, nil)
	logrus.Fatal(server.ListenAndServe())
}
