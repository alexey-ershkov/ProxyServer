package handlers

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func doHttpProxyRequest(request *http.Request) (*http.Response, error) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	proxyReq, err := http.NewRequest(request.Method, request.RequestURI, request.Body)
	if err != nil {
		return nil, err
	}

	for reqHeaderKey, reqHeaderValues := range request.Header {
		for _, reqHeaderValue := range reqHeaderValues {
			if !contains(skipHeaderList, strings.ToLower(reqHeaderKey)) {
				proxyReq.Header.Add(reqHeaderKey, reqHeaderValue)
			}
		}
	}

	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		return nil, err
	}

	return proxyResp, nil
}

func sendHttpProxyResponse (responseWriter http.ResponseWriter,response *http.Response) {
	responseWriter.WriteHeader(response.StatusCode)
	copyHeaders(response.Header, responseWriter.Header())

	_, err := io.Copy(responseWriter, response.Body)
	if err != nil {
		logrus.Error(err)
	}
}
