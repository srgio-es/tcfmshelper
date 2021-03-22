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

func parseStatus(status string) model.FscStatus {
	var fscStatus model.FscStatus

	linesSplited := cleanAndSplitOutput(status)

	switch {
	case linesSplited[3] == "true":
		fscStatus.Status = model.STATUS_OK
		return fscStatus

	default:
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

}

func tailLog(output string, lines int64) string {
	var result string

	splitted := cleanAndSplitOutput(output)
	tail := splitted[len(splitted)-int(lines):]

	result += fmt.Sprintln(splitted[3] + "\n")

	for _, s := range tail {
		result += fmt.Sprintln(s)
	}

	return result
}

func parseVersion(output string) model.FSCVersion {
	var fscVersion model.FSCVersion

	linesSplited := cleanAndSplitOutput(output)

	fscVersion.FmsServerCache.Version = linesSplited[3][strings.Index(linesSplited[3], ":")+2 : strings.Index(linesSplited[3], ",")]
	fscVersion.FmsServerCache.BuildDate = linesSplited[3][strings.LastIndex(linesSplited[3], ":")+2:]

	fscVersion.FmsUtil.Version = linesSplited[4][strings.Index(linesSplited[4], ":")+2 : strings.Index(linesSplited[4], ",")]
	fscVersion.FmsUtil.BuildDate = linesSplited[4][strings.LastIndex(linesSplited[4], ":")+2:]

	fscVersion.FscJavaClientProxy.Version = linesSplited[5][strings.Index(linesSplited[5], ":")+2 : strings.Index(linesSplited[5], ",")]
	fscVersion.FscJavaClientProxy.BuildDate = linesSplited[5][strings.LastIndex(linesSplited[5], ":")+2:]

	return fscVersion
}

func parseHash(output string) string {
	return cleanAndSplitOutput(output)[3]
}

func parseConfigReport(output string) []model.FscConfig {
	var result []model.FscConfig

	linesSplited := cleanAndSplitOutput(output)

	mastersIndex := indexOf("# masters", linesSplited)
	slavesindex := indexOf("# slaves", linesSplited)

	masters := linesSplited[mastersIndex+1 : slavesindex]
	slaves := linesSplited[slavesindex+1:]

	for _, item := range masters {
		result = append(result, convertToConfigItem(item))
	}

	for _, item := range slaves {
		result = append(result, convertToConfigItem(item))
	}

	return result
}

func convertToConfigItem(configline string) model.FscConfig {
	var item model.FscConfig
	var err error

	splitted := strings.Split(configline, ",")

	item.FSCId = splitted[0]
	item.ConfigHash = splitted[1]

	item.IsMaster, err = strconv.ParseBool(splitted[2])
	if err != nil {
		log.Printf("Failed while parsing FSC config item: %v", err)
	}

	if splitted[3] == "ok" {
		item.Status = model.STATUS_OK
	} else {
		item.Status = model.STATUS_KO
		//The following is not the prettiest but is done to keep output split straightforward.
		errors := splitted[3:]
		for _, e := range errors {
			item.Error += "," + e
		}
	}

	return item
}

func parseError(output string) error {
	var err error

	// lines := strings.ReplaceAll(output, "\n", "")
	linesSplited := strings.Split(output, "\n")

	nativeError := strings.ReplaceAll(linesSplited[3], "\r", "")

	switch {
	case strings.Contains(nativeError, "java.net.UnknownHostException"):
		unkownHost := nativeError[strings.Index(nativeError, ":")+2:]
		err = fmt.Errorf("Unknown Host: %s", unkownHost)
	case strings.Contains(nativeError, "java.net.ConnectException"):
		unkownHost := nativeError[strings.Index(nativeError, ":")+2:]
		err = fmt.Errorf("%s", unkownHost)
	case strings.Contains(nativeError, "java.net.MalformedURLException"):
		malformedUri := nativeError[strings.Index(nativeError, ":")+2:]
		err = fmt.Errorf("Malformed URI: %s", malformedUri)
	}

	return err
}

func cleanAndSplitOutput(output string) []string {
	lines := strings.ReplaceAll(output, "\n", "")
	linesSplited := strings.Split(lines, "\r")

	if len(linesSplited) > 0 {
		linesSplited = linesSplited[:len(linesSplited)-1]
	}

	return linesSplited
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
