package handlers

import (
	"math/rand"
	"os/exec"
	"strconv"
)

func genProxyCert(scriptPath, scriptName, host, savePath string) error {
	genCmd := exec.Command(scriptPath+scriptName, host, strconv.Itoa(rand.Intn(100000000)), scriptPath, savePath)
	_, err := genCmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
