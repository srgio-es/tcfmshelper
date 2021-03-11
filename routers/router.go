package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/srgio-es/tcfmshelper/fscadmin"
	"github.com/srgio-es/tcfmshelper/settings"
)

//InitRouter initializes main router
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/serverhealth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"server_status": "up",
		})
	})

	r.GET("/fscstatus", func(c *gin.Context) {

		fscCommand := fscadmin.FscCommand{
			JavaHome:   settings.AppSettings.JavaHome,
			FmsHome:    settings.FscSettings.FscLocation,
			FscFromUrl: settings.FscSettings.FscFromUrl,
		}

		status := fscCommand.FSCStatus()

		c.JSON(200, status)
	})

	return r
}
