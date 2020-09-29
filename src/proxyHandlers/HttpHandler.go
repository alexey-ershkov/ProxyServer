package proxyHandlers

import (
	"Proxy/src/db"
	"Proxy/src/models"
	"Proxy/src/utils"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

type HttpHandler struct {
	respWriter    http.ResponseWriter
	clientRequest *http.Request
	proxyResp     *http.Response
	dbConn        *db.Database
}

func NewHttpHandler(respWriter http.ResponseWriter, clientRequest *http.Request, dbConn *db.Database) *HttpHandler {
	return &HttpHandler{
		respWriter:    respWriter,
		clientRequest: clientRequest,
		dbConn:        dbConn,
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

func (hh *HttpHandler) Defer() {
	hh.dbConn.Close()
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

	reqDump, err := httputil.DumpRequest(proxyReq, true)
	if err != nil {
		logrus.Warn("Can't dump request")
	}

	dbReq := models.DatabaseReq{
		Host:    hh.clientRequest.RequestURI,
		Request: string(reqDump),
	}
	hh.dbConn.InsertRequest(dbReq)

	hh.proxyResp, err = client.Do(proxyReq)
	if err != nil {
		return err
	}

	return nil
}

func (hh *HttpHandler) sendResponse() error {
	utils.CopyHeaders(hh.proxyResp.Header, hh.respWriter.Header())
	hh.respWriter.WriteHeader(hh.proxyResp.StatusCode)
	hh.respWriter.Header().Add("Connection", hh.proxyResp.Header.Get("Connection"))

	_, err := io.Copy(hh.respWriter, hh.proxyResp.Body)
	if err != nil {
		return err
	}
	return nil
}
