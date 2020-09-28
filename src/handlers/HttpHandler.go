package handlers

import (
	"io"
	"net/http"
	"strings"
)

type HttpHandler struct {
	respWriter    http.ResponseWriter
	clientRequest *http.Request
	proxyResp     *http.Response
}

func NewHttpHandler(respWriter http.ResponseWriter, clientRequest *http.Request) *HttpHandler {
	return &HttpHandler{
		respWriter:    respWriter,
		clientRequest: clientRequest,
	}
}

func (hh *HttpHandler) ProxyRequest() error {
	err := hh.doRequest()
	if err != nil {
		return err
	}

	err = hh.sendResponse()
	if err != nil {
		return err
	}

	return nil
}

func (hh *HttpHandler) Defer () {
}

func (hh *HttpHandler) doRequest() error {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	proxyReq, err := http.NewRequest(hh.clientRequest.Method, hh.clientRequest.RequestURI, hh.clientRequest.Body)
	if err != nil {
		return err
	}

	for reqHeaderKey, reqHeaderValues := range hh.clientRequest.Header {
		for _, reqHeaderValue := range reqHeaderValues {
			if !contains(skipHeaderList, strings.ToLower(reqHeaderKey)) {
				proxyReq.Header.Add(reqHeaderKey, reqHeaderValue)
			}
		}
	}

	hh.proxyResp, err = client.Do(proxyReq)
	if err != nil {
		return err
	}

	return nil
}

func (hh *HttpHandler) sendResponse() error {
	copyHeaders(hh.proxyResp.Header, hh.respWriter.Header())
	hh.respWriter.WriteHeader(hh.proxyResp.StatusCode)
	hh.respWriter.Header().Add("Connection", hh.proxyResp.Header.Get("Connection"))

	_, err := io.Copy(hh.respWriter, hh.proxyResp.Body)
	if err != nil {
		return err
	}
	return nil
}
