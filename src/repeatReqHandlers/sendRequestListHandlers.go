package repeatReqHandlers

import (
	"Proxy/db"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SendRequestList(respWriter http.ResponseWriter, _ *http.Request) {
	dbConn, err := db.CreateNewDatabaseConnection()
	if err != nil {
		logrus.Warn("Can't connect to database")
		logrus.Error(err)
	}

	defer dbConn.Close()

	reqList, err := dbConn.GetRequestList()
	if err != nil {
		logrus.Warn("Can't get data from DB")
		_, _ = fmt.Fprintf(respWriter, "Can't get request info\n")
		return
	}

	if len(reqList) == 0 {
		_, _ = fmt.Fprintf(respWriter, "No requests saved\n")
		return
	}

	_, _ = fmt.Fprintf(respWriter,
		"Saved requests\n--------------\n",
	)
	for i, req := range reqList {
		_, _ = fmt.Fprintf(respWriter,
			"%d) Host: %s\n\n%s"+
				"------ RepeatLink: http://127.0.0.1/req?id=%d\n"+
				"------ Example: curl -X GET \"http://127.0.0.1/req?id=%d\"\n\n",
			i+1,
			req.Host,
			req.Request,
			req.Id,
			req.Id)
	}
}
