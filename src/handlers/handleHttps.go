package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

func handleHTTPS(request *http.Request) (*http.Response, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error(err)
	}

	genScriptDir := pwd + "/certGen"
	certsDir := pwd + "/certs/"

	parsedUrl, err := url.Parse(request.RequestURI)
	if err != nil {
		logrus.Error(err)
	}

	genCmd := exec.Command(genScriptDir+"/gen_cert.sh", parsedUrl.Scheme, "01", genScriptDir, certsDir)
	info, err := genCmd.CombinedOutput()
	if err != nil {
		logrus.Error(string(info), err)
	} else {
		logrus.Debug(string(info))
	}
	return nil, nil
}
