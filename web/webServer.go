package web

import (
	// "strconv"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.GET("/Process", func(c *gin.Context) {
		// 获取下载和延迟的当前值和总值
		currentDownload, totalDownload := GetProcessDownloadBar()
		currentDelay, totalDelay := GetProcessDelayBar()

		// 返回 JSON 格式的数据
		c.JSON(200, gin.H{
			"Download": gin.H{
				"Current": currentDownload,
				"Total":   totalDownload,
			},
			"Delay": gin.H{
				"Current": currentDelay,
				"Total":   totalDelay,
			},
		})
	})

	r.GET("/AllData", func(c *gin.Context) {
		c.JSON(200, GetAllData())
	})

	r.GET("/Schedules", func(c *gin.Context) {
		times := GetSchedules()
		c.JSON(200, times)
	})

	r.Run() // 启动服务，监听和服务在 0.0.0.0:8080
}
