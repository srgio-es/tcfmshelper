package routers

import (
	"log"
	"net/http"
	"strings"

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

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	fscCommand := fscadmin.FscCommand{
		JavaHome: settings.AppSettings.JavaHome,
		FmsHome:  settings.FscSettings.FscLocation,
	}

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"server_status": "up",
		})
	})

	r.GET("/fscstatus/:host", func(c *gin.Context) {

		host := c.Param("host")
		port := c.DefaultQuery("port", "4544")

		status, err := fscCommand.FSCStatus(host, port)

		if err != nil {
			c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
		} else {
			c.JSON(200, status)
		}
	})

	r.GET("/fscalive/:host", func(c *gin.Context) {
		host := c.Param("host")
		port := c.DefaultQuery("port", "4544")

		status, err := fscCommand.FCSAlive(host, port)

		if err != nil {
			c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
		} else {
			c.JSON(200, status)
		}
	})

	r.GET("/fscversion/:host", func(c *gin.Context) {
		host := c.Param("host")
		port := c.DefaultQuery("port", "4544")

		status, err := fscCommand.FSCVersion(host, port)

		if err != nil {
			c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
		} else {
			c.JSON(200, status)
		}
	})

	r.GET("/fscconfig/:host", func(c *gin.Context) {
		host := c.Param("host")
		port := c.DefaultQuery("port", "4544")

		output, err := fscCommand.FSCConfig(host, port)

		if err != nil {
			c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
		} else {
			// c.String(200, output)
			c.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
			c.String(200, output)
		}
	})

	r.GET("/fscconfig/:host/hash/", func(c *gin.Context) {
		host := c.Param("host")
		port := c.DefaultQuery("port", "4544")

		output, err := fscCommand.FSCConfigHash(host, port)

		if err != nil {
			c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
		} else {
			c.String(200, output)
		}
	})

	r.GET("/fscconfig/", func(c *gin.Context) {
		for i, fsc := range settings.FscSettings.FmsMasterURL {
			host := fsc[:strings.Index(fsc, ":")]
			port := fsc[strings.Index(fsc, ":")+1:]

			log.Printf("host: %s", host)
			log.Printf("port: %s", port)

			output, err := fscCommand.FSCConfig(host, port)
			if err != nil && !strings.Contains(err.Error(), "Unknown Host") {
				c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
			} else if err != nil && i == len(settings.FscSettings.FmsMasterURL)-1 {
				c.JSON(500, model.Error{Status: model.STATUS_KO, Message: "All declared FMS masters are down"})
			} else if err == nil {
				c.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
				c.String(200, output)
				return
			}
		}
	})

	r.GET("/fmsconfigreport/", func(c *gin.Context) {
		for i, fsc := range settings.FscSettings.FmsMasterURL {
			host := fsc[:strings.Index(fsc, ":")]
			port := fsc[strings.Index(fsc, ":")+1:]

			log.Printf("host: %s", host)
			log.Printf("port: %s", port)

			output, err := fscCommand.FSCConfigReport(host, port)
			if err != nil && !strings.Contains(err.Error(), "Unknown Host") {
				c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
			} else if err != nil && i == len(settings.FscSettings.FmsMasterURL)-1 {
				c.JSON(500, model.Error{Status: model.STATUS_KO, Message: "All declared FMS masters are down"})
			} else if err == nil {
				c.JSON(200, output)
				return
			}
		}
	})

	r.GET("/fscstatus/", func(c *gin.Context) {

		for i, fsc := range settings.FscSettings.FmsMasterURL {
			host := fsc[:strings.Index(fsc, ":")]
			port := fsc[strings.Index(fsc, ":")+1:]

			log.Printf("host: %s", host)
			log.Printf("port: %s", port)

			output, err := fscCommand.FSCStatusAll(host, port, settings.FscSettings.MaxParallel)
			if err != nil && !strings.Contains(err.Error(), "Unknown Host") {
				c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
			} else if err != nil && i == len(settings.FscSettings.FmsMasterURL)-1 {
				c.JSON(500, model.Error{Status: model.STATUS_KO, Message: "All declared FMS masters are down"})
			} else if err == nil {
				c.JSON(200, output)
				return
			}
		}

	})

	r.GET("/fscdashboard/", func(c *gin.Context) {
		for i, fsc := range settings.FscSettings.FmsMasterURL {
			host := fsc[:strings.Index(fsc, ":")]
			port := fsc[strings.Index(fsc, ":")+1:]

			log.Printf("host: %s", host)
			log.Printf("port: %s", port)

			output, err := fscCommand.FSCStatusAll(host, port, settings.FscSettings.MaxParallel)
			if err != nil && !strings.Contains(err.Error(), "Unknown Host") {
				c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
			} else if err != nil && i == len(settings.FscSettings.FmsMasterURL)-1 {
				c.JSON(500, model.Error{Status: model.STATUS_KO, Message: "All declared FMS masters are down"})
			} else if err == nil {
				c.HTML(http.StatusOK, "fscdashboard.tmpl", output)
				return
			}
		}
	})

	r.GET("/log/:host", func(c *gin.Context) {
		host := c.Param("host")
		port := c.DefaultQuery("port", "4544")

		lines := c.DefaultQuery("lines", "all")

		output, err := fscCommand.FSCLog(host, port, lines)

		if err != nil {
			c.JSON(500, model.Error{Status: model.STATUS_KO, Message: err.Error()})
		} else {
			// c.String(200, output)
			c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
			c.String(200, output)
		}
	})

	return r
}
