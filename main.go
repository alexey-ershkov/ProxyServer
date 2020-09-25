package main

import "net/http"
import "github.com/sirupsen/logrus"


func mainHandler(writer http.ResponseWriter, reader *http.Request){
	logrus.Debug("Handle",reader.Header, reader.Host, reader.GetBody)
	writer.WriteHeader(200)
	reader.Write(writer)
}

func main()  {
	logrus.SetLevel(logrus.InfoLevel)

	server := &http.Server{
		Addr: ":3000",
		Handler: http.HandlerFunc(mainHandler),
	}

	logrus.Info("Server started on port 3000")

	logrus.Fatal(server.ListenAndServe())
}
