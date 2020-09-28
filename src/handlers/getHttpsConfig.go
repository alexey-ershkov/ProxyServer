package handlers

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func getHttpsConfig(host string) (*tls.Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genScriptAndRootCaDir := pwd + "/certGen"
	certsDir := pwd + "/certs/"

	err = genProxyCert(genScriptAndRootCaDir, "/gen_cert.sh", host, certsDir)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(certsDir+host+".crt", genScriptAndRootCaDir+"/cert.key")
	if err != nil {
		logrus.Error(err)
	}

	config := new(tls.Config)
	config.Certificates = []tls.Certificate{cert}
	config.ServerName = host

	return config, nil
}

func genProxyCert(scriptPath, scriptName, host, savePath string) error {
	genCmd := exec.Command(scriptPath+scriptName, host, "01", scriptPath, savePath)
	info, err := genCmd.CombinedOutput()
	if err != nil {
		return err
	}
	logrus.Debug(info)
	return nil
}
