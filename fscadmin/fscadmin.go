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
	JavaHome   string
	FmsHome    string
	FscFromUrl string
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

	if err != nil {
		log.Printf("Error executing fscadmin: %v", err)
		return "", err
	}

	return fmt.Sprintf("%s", output), err
}

func (fsc *FscCommand) FSCStatus() model.FscStatus {
	result, err := fsc.fscAdminExec("-s", fsc.FscFromUrl, "./status")
	if err != nil {
		log.Printf("Error executing fscadmin: %v", err)
	}

	return parseStatus(result)

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

	log.Printf("%#v", fscStatus)

	return fscStatus
}
