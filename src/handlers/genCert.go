package handlers

import (
	"os/exec"
)

func genProxyCert(scriptPath, scriptName, host, savePath string) error {
	genCmd := exec.Command(scriptPath+scriptName, host, "01", scriptPath, savePath)
	_, err := genCmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
