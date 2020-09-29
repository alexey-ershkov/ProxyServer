package repeatReqHandlers

import (
	"Proxy/db"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	fmt.Println(req)

}
