package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/srgio-es/tcfmshelper/routers"
	"github.com/srgio-es/tcfmshelper/settings"
)

func Run() error {
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

	err := server.ListenAndServe()

	return err
}
