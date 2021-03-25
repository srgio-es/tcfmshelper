package main

import (
	"log"
	"os"

	"github.com/srgio-es/tcfmshelper/server"
	"github.com/srgio-es/tcfmshelper/settings"
)

func init() {
	settings.Setup()
}

func main() {
	err := server.Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
