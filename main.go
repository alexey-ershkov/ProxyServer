package main

import (
	"Proxy/src/proxyHandlers"
	"Proxy/src/repeatReqHandlers"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func init() {
	debugLevel, _ := os.LookupEnv("DEBUG_LEVEL")
	switch debugLevel {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	if err := godotenv.Load(); err != nil {
		logrus.Print("No .env file found")
	}
}

func main() {

	ProxyPort, _ := os.LookupEnv("PROXY_PORT")
	RepeatPort, _ := os.LookupEnv("REPEAT_PORT")

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
