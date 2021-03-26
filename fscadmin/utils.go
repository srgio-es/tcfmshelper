package fscadmin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/srgio-es/tcfmshelper/fscadmin/model"
	"github.com/srgio-es/tcfmshelper/settings"
	"go.uber.org/zap"
)

func parseStatus(status string, host string) model.FscStatus {
	var fscStatus model.FscStatus

	linesSplited := cleanAndSplitOutput(status)

	fscStatus.Host = host

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
			settings.Log.Logger.Error("Failed while parsing FSC Status", zap.Error(err))
		}
		fscStatus.CurrentAdminConnections, err = strconv.ParseInt(linesSplited[6][strings.Index(linesSplited[6], ":")+2:], 10, 64)
		if err != nil {
			settings.Log.Logger.Error("Failed while parsing FSC Status", zap.Error(err))

		}

		//TODO: Parse duration

		fscStatus.Status = model.STATUS_OK

		settings.Log.Logger.Debug("FscStatus", zap.Any("value", fscStatus))

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
		settings.Log.Logger.Error("Failed while parsing FSC config item", zap.Error(err))
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
