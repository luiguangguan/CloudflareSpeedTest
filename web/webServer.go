package web

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.GET("/Process", func(c *gin.Context) {
		current := GetProcessDownloadBar()
		c.String(200, strconv.FormatInt(current, 10))
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
