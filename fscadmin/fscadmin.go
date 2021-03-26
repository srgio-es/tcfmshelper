package fscadmin

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/antchfx/xmlquery"
	"github.com/srgio-es/tcfmshelper/fscadmin/model"
	"github.com/srgio-es/tcfmshelper/settings"
	"go.uber.org/zap"
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

	settings.Log.Logger.Debug("FSCAdmin CMD input", zap.Strings("args", cmd.Args))

	output, err := cmd.CombinedOutput()

	return fmt.Sprintf("%s", output), err
}

func (fsc *FscCommand) FSCStatus(host string, port string) (model.FscStatus, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./status")
	if err != nil {
		settings.Log.Logger.Error("Error executing FSCStatus", zap.Error(err), zap.String("result", output))
		return model.FscStatus{Host: host}, parseError(output)
	}

	return parseStatus(output, host), nil

}

func (fsc *FscCommand) FSCStatusAll(host string, port string, parallelWorkers int) ([]model.FscStatus, error) {

	configXml, err := fsc.FSCConfig(host, port)
	if err != nil {
		return nil, err
	}

	doc, err := xmlquery.Parse(strings.NewReader(configXml))
	if err != nil {
		settings.Log.Logger.Error("Error while parsing FSC Config XML", zap.Error(err))
		return nil, err
	}

	nodes := make(chan *xmlquery.Node)

	go func() {
		defer close(nodes)
		for _, node := range xmlquery.Find(doc, "//fsc") {
			nodes <- node
		}

	}()

	// fscNodes := xmlquery.Find(doc, "//fsc")

	statuses := make(chan model.FscStatus)

	nWorkers := parallelWorkers
	workers := int32(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go func() {
			defer func() {
				// Last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					close(statuses)
				}
			}()

			for n := range nodes {
				fmt.Printf("%v\n", n.SelectAttr("address"))
				addr := n.SelectAttr("address")
				h := addr[strings.Index(addr, "//")+2 : strings.LastIndex(addr, ":")]
				p := addr[strings.LastIndex(addr, ":")+1:]
				settings.Log.Logger.Debug("FSCStatusAll Node called", zap.String("host", h), zap.String("port", p))

				status, err := fsc.FSCStatus(h, p)
				if err != nil {
					status = model.FscStatus{Host: h, Status: model.STATUS_KO, Error: err.Error()}
				}

				statuses <- status
			}
		}()
	}

	ret := []model.FscStatus{}
	for status := range statuses {
		ret = append(ret, status)
	}

	return ret, nil
}

func (fsc *FscCommand) FSCConfig(host string, port string) (string, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./config")
	if err != nil {
		settings.Log.Logger.Error("Error executing FSCConfig", zap.Error(err), zap.String("result", output))
		return "", parseError(output)
	}

	return output[strings.Index(output, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") : strings.LastIndex(output, "</fmsworld>")+len("</fmsworld>")], nil
}

func (fsc *FscCommand) FSCLog(host string, port string, lines string) (string, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./log")
	if err != nil {
		settings.Log.Logger.Error("Error executing FSCLog", zap.Error(err), zap.String("result", output))
		return "", parseError(output)
	}

	if lines == "all" {
		return output[strings.Index(output, "Current log file:"):], nil
	} else {
		li, err := strconv.ParseInt(lines, 10, 64)
		if err != nil {
			settings.Log.Logger.Error("Error executing FSCLog", zap.Error(err))
			return "", errors.New("Parameter lines has to be 'all' or a valid integer")
		}
		return tailLog(output, li), nil
	}

}

func (fsc *FscCommand) FCSAlive(host string, port string) (model.FscStatus, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./alive")
	if err != nil {
		settings.Log.Logger.Error("Error executing FSCAlive", zap.Error(err), zap.String("result", output))
		return model.FscStatus{}, parseError(output)
	}

	return parseStatus(output, host), nil
}

func (fsc *FscCommand) FSCVersion(host string, port string) (model.FSCVersion, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./version")
	if err != nil {
		settings.Log.Logger.Error("Error executing FscVersion", zap.Error(err), zap.String("result", output))
		return model.FSCVersion{}, parseError(output)
	}

	return parseVersion(output), nil
}

func (fsc *FscCommand) FSCConfigHash(host string, port string) (string, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./config/hash")
	if err != nil {
		settings.Log.Logger.Error("Error executing FscConfigHash", zap.Error(err), zap.String("result", output))
		return "", parseError(output)
	}
	return parseHash(output), nil
}

func (fsc *FscCommand) FSCConfigReport(host string, port string) ([]model.FscConfig, error) {
	url := "http://" + host + ":" + port
	output, err := fsc.fscAdminExec("-s", url, "./config/report")
	if err != nil {
		settings.Log.Logger.Error("Error executing FscConfigReport", zap.Error(err), zap.String("result", output))
		return nil, parseError(output)
	}
	return parseConfigReport(output), nil
}
