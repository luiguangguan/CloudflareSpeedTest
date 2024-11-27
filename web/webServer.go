package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var writeMutex sync.Mutex

// 连接管理
var connections = make(map[*websocket.Conn]bool)
var connMutex = sync.Mutex{}

func Start() {
	r := gin.Default()

	// 提供 Vue 构建后的静态文件
	r.Static("/static", "./static/vue")
	r.Static("/assets", "./static/vue/assets")

	// 根路径访问 index.html
	r.GET("/speed", func(c *gin.Context) {
		c.File("./static/vue/index.html")
	})

	// WebSocket 接口
	r.GET("/Process", func(c *gin.Context) {
		// 升级连接为 WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
			return
		}

		// 添加连接到管理列表
		connMutex.Lock()
		connections[conn] = true
		connMutex.Unlock()

		fmt.Println("New connection established.")

		// 启动数据推送
		go handleProcessConnection(conn)
	})

	// 其他路由保留不变
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

	r.GET("/Get1DayMaxData", func(c *gin.Context) {
		c.JSON(200, Get1DayMaxData())
	})
	r.GET("/Get3DayMaxData", func(c *gin.Context) {
		c.JSON(200, Get3DayMaxData())
	})
	r.GET("/Get5DayMaxData", func(c *gin.Context) {
		c.JSON(200, Get5DayMaxData())
	})

	r.GET("/GetYesterdayMaxData", func(c *gin.Context) {
		c.JSON(200, GetYesterdayMaxData())
	})

	// 启动服务
	r.Run(":8080")
}

func handleProcessConnection(conn *websocket.Conn) {
	defer func() {
		// 移除连接并关闭
		connMutex.Lock()
		delete(connections, conn)
		connMutex.Unlock()
		conn.Close()
		fmt.Println("Connection closed and removed.")
	}()

	// 定时器：心跳发送
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() { // 心跳检测线程
		for range ticker.C {
			writeMutex.Lock()
			if err := conn.WriteMessage(websocket.TextMessage, []byte("heartbeat")); err != nil {
				fmt.Printf("Heartbeat failed: %v\n", err)
				writeMutex.Unlock()
				conn.Close()
				return
			}
			writeMutex.Unlock()
		}
	}()

	// 数据推送逻辑
	for {
		// 监听客户端是否断开
		// _, _, err := conn.ReadMessage()
		// if err != nil {
		// 	fmt.Printf("Client disconnected: %v\n", err)
		// 	break
		// }

		// 模拟获取数据并推送
		currentDownload, totalDownload, downloadIP, speed := GetProcessDownloadBar()
		currentDelay, totalDelay, delayIP, available := GetProcessDelayBar()
		count := GetAllDataCount()

		message := gin.H{
			"Download": gin.H{
				"Current": currentDownload,
				"Total":   totalDownload,
				"IP":      downloadIP,
				"Speed":   speed,
			},
			"Delay": gin.H{
				"Current":   currentDelay,
				"Total":     totalDelay,
				"IP":        delayIP,
				"Available": available,
			},
			"AllDataCount": count,
		}

		// 写入数据
		writeMutex.Lock()
		if err := conn.WriteJSON(message); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			writeMutex.Unlock()
			break
		}
		writeMutex.Unlock()

		// 控制推送速度
		if currentDownload == totalDownload && currentDelay == totalDelay {
			time.Sleep(1 * time.Second) // 推送间隔
		} else {
			time.Sleep(300 * time.Millisecond)
		}
	}
}
