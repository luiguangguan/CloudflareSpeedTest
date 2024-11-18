package web

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.GET("/Process", func(c *gin.Context) {
		currentDownload := GetProcessDownloadBar()
		currentDelay := GetProcessDelayBar()

		c.String(200, `{"Download":`+strconv.FormatInt(currentDownload, 10)+`,"Delay":`+strconv.FormatInt(currentDelay, 10)+`}`)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
