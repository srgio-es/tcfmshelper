package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/srgio-es/tcfmshelper/routers"
	"github.com/srgio-es/tcfmshelper/settings"
)

func init() {
	log.Println("Starting tcfmshelper")
	settings.Setup()
}

func main() {
	gin.SetMode(settings.ServerSettings.RunMode)

	routers := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", settings.ServerSettings.Port)

	server := &http.Server{
		Addr:         endPoint,
		Handler:      routers,
		ReadTimeout:  settings.ServerSettings.ReadTimeout,
		WriteTimeout: settings.ServerSettings.WriteTimeout,
	}

	log.Printf("Start http server listening %s", endPoint)

	server.ListenAndServe()
}
