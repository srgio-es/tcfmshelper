package fscadmin

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/srgio-es/tcfmshelper/fscadmin/model"
)

type FscCommand struct {
	JavaHome string
	FmsHome  string
}

func (fsc *FscCommand) fscAdminExec(args ...string) (string, error) {

	javaExecutable := filepath.Join(fsc.JavaHome, "bin", "java")

	fmsLibs := filepath.Join(fsc.FmsHome, "jar", "fmsutil.jar") + ";" + filepath.Join(fsc.FmsHome, "jar", "fmsservercache.jar")

	cmd := exec.Command(javaExecutable, "-classpath", fmsLibs, "com.teamcenter.fms.servercache.fscadmin.FSCAdmin")

	for _, arg := range args {
		cmd.Args = append(cmd.Args, arg)
	}

	log.Println(cmd.Args)

	output, err := cmd.CombinedOutput()

	return fmt.Sprintf("%s", output), err
}

func (fsc *FscCommand) FSCStatus(host string) (model.FscStatus, error) {
	url := "http://" + host
	output, err := fsc.fscAdminExec("-s", url, "./status")
	if err != nil {
		log.Printf("Error executing FSCStatus error: %v\nresult:%v", err, output)
		return model.FscStatus{}, parseError(output)
	}

	return parseStatus(output), nil

}

func parseStatus(status string) model.FscStatus {
	var fscStatus model.FscStatus

	lines := strings.ReplaceAll(status, "\n", "")
	linesSplited := strings.Split(lines, "\r")

	log.Printf("%#v", linesSplited)

	fscStatus.FSCId = linesSplited[3][8:strings.Index(linesSplited[3], ",")]
	fscStatus.Site = linesSplited[3][strings.Index(linesSplited[3], ",")+7:]

	var err error
	fscStatus.CurrentFileConnections, err = strconv.ParseInt(linesSplited[5][strings.Index(linesSplited[5], ":")+2:], 10, 64)
	if err != nil {
		log.Printf("Failed while parsing FSC Status result: %v", err)
	}
	fscStatus.CurrentAdminConnections, err = strconv.ParseInt(linesSplited[6][strings.Index(linesSplited[6], ":")+2:], 10, 64)
	if err != nil {
		log.Printf("Failed while parsing FSC Status result: %v", err)
	}

	//TODO: Parse duration

	fscStatus.Status = model.STATUS_OK

	log.Printf("%#v", fscStatus)

	return fscStatus
}

func parseError(output string) error {
	var err error

	lines := strings.ReplaceAll(output, "\n", "")
	linesSplited := strings.Split(lines, "\r")

	nativeError := linesSplited[3]

	if strings.Contains(nativeError, "java.net.UnknownHostException") {
		unkownHost := nativeError[strings.Index(nativeError, ":")+1:]
		err = fmt.Errorf("Unknown Host: %s", unkownHost)
	}

	return err
}
