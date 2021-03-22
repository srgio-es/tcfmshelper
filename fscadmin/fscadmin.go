package fscadmin

import (
	"errors"
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

func (fsc *FscCommand) FSCStatus(host string, port string) (model.FscStatus, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./status")
	if err != nil {
		log.Printf("Error executing FSCStatus error: %v\nresult: %v", err, output)
		return model.FscStatus{}, parseError(output)
	}

	return parseStatus(output), nil

}

func (fsc *FscCommand) FSCConfig(host string, port string) (string, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./config")
	if err != nil {
		log.Printf("Error executing FSCConfig error: %v\nresult: %v", err, output)
		return "", parseError(output)
	}

	return output[strings.Index(output, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") : strings.LastIndex(output, "</fmsworld>")+len("</fmsworld>")], nil
}

func (fsc *FscCommand) FSCLog(host string, port string, lines string) (string, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./log")
	if err != nil {
		log.Printf("Error executing FSCLog error: %v\nresult: %v", err, output)
		return "", parseError(output)
	}

	if lines == "all" {
		return output[strings.Index(output, "Current log file:"):], nil
	} else {
		li, err := strconv.ParseInt(lines, 10, 64)
		if err != nil {
			log.Printf("Error executing FSCLog error: %v\n", err)
			return "", errors.New("Parameter lines has to be 'all' or a valid integer")
		}
		return tailLog(output, li), nil
	}

}

func (fsc *FscCommand) FCSAlive(host string, port string) (model.FscStatus, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./alive")
	if err != nil {
		log.Printf("Error executing FSCAlive error: %s\nresult: %v", err, output)
		return model.FscStatus{}, parseError(output)
	}

	return parseStatus(output), nil
}

func (fsc *FscCommand) FSCVersion(host string, port string) (model.FSCVersion, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./version")
	if err != nil {
		log.Printf("Error executing FscVersion error: %s\nresult: %v", err, output)
		return model.FSCVersion{}, parseError(output)
	}

	return parseVersion(output), nil
}

func (fsc *FscCommand) FSCConfigHash(host string, port string) (string, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./config/hash")
	if err != nil {
		log.Printf("Error executing FscVersion error: %s\nresult: %v", err, output)
		return "", parseError(output)
	}
	return parseHash(output), nil
}

func (fsc *FscCommand) FSCConfigReport(host string, port string) ([]model.FscConfig, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./config/report")
	if err != nil {
		log.Printf("Error executing FscVersion error: %s\nresult: %v", err, output)
		return nil, parseError(output)
	}
	return parseConfigReport(output), nil
}
