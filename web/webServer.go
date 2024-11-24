package web

import (
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	// 提供 Vue 构建后的静态文件
	r.Static("/static", "./static/vue")        // 提供 Vue 构建后的静态文件
	r.Static("/assets", "./static/vue/assets") // 提供 Vue 构建后的 assets 文件

	// 让根路径（/）访问 index.html
	r.GET("/", func(c *gin.Context) {
		c.File("./static/vue/index.html") // 直接返回 Vue 构建后的 index.html
	})

	// 后端 API 接口
	r.GET("/Process", func(c *gin.Context) {
		currentDownload, totalDownload := GetProcessDownloadBar()
		currentDelay, totalDelay := GetProcessDelayBar()

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

	r.GET("/MaxData", func(c *gin.Context) {
		c.JSON(200, GetMaxData())
	})

	// 启动服务，监听和服务在 0.0.0.0:8080
	r.Run(":8080")
}
