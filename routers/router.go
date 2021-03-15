package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/srgio-es/tcfmshelper/fscadmin"
	"github.com/srgio-es/tcfmshelper/fscadmin/model"
	"github.com/srgio-es/tcfmshelper/settings"
)

var fscCommand fscadmin.FscCommand

//InitRouter initializes main router
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	fscCommand := fscadmin.FscCommand{
		JavaHome: settings.AppSettings.JavaHome,
		FmsHome:  settings.FscSettings.FscLocation,
	}

	r.GET("/serverhealth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"server_status": "up",
		})
	})

	r.GET("/fscstatus/:host", func(c *gin.Context) {

		host := c.Param("host")

		status, err := fscCommand.FSCStatus(host)

		if err != nil {
			c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
		} else {
			c.JSON(200, status)
		}
	})

	r.GET("/fscalive/:host", func(c *gin.Context) {
		host := c.Param("host")

		status, err := fscCommand.FCSAlive(host)

		if err != nil {
			c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
		} else {
			c.JSON(200, status)
		}
	})

	return r
}
