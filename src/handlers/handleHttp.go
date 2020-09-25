package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func handleHTTP(request *http.Request) *http.Response{
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	proxyReq, err := http.NewRequest(request.Method, request.RequestURI, request.Body)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
	}

	return proxyResp
}
