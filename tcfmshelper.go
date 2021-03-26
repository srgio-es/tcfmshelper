package main

import (
	"os"

	"github.com/srgio-es/tcfmshelper/server"
	"github.com/srgio-es/tcfmshelper/settings"
	"go.uber.org/zap"
)

func init() {
	settings.Setup()
}

func main() {
	settings.Log.Logger.Sugar().Info("TCFMSHELPER starting")
	err := server.Run()
	if err != nil {
		settings.Log.Logger.Error("An error occurred", zap.Error(err))
		os.Exit(1)
	}
}
