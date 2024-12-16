package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

	r.GET("/TraceInfosCount", func(c *gin.Context) {
		d1, alldata := TraceInfosCount()
		c.JSON(200, gin.H{
			"AllData": alldata,
			"Count":   d1,
		})
	})

	r.GET("/GetIPTraceInfos", func(c *gin.Context) {
		c.JSON(200, GetIPTraceInfos())
	})

	r.POST("/ClearTracerInfo", func(c *gin.Context) {
		var data Pwd

		// 绑定 JSON 请求体到 SubmitData 结构体
		if err := c.ShouldBindJSON(&data); err != nil {
			// 如果请求体有问题，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request data",
			})
			return
		}
		b, _ := utils.CheckPassword(data.Password)
		if b {
			ok, e := ClearTracerInfo()
			if ok {
				c.JSON(200, gin.H{
					"message": "删除成功",
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "失败" + e,
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "密码错误",
			})
		}
	})

	// r.GET("/GetPwd", func(c *gin.Context) {
	// 	d, _ := utils.Select("select * from Passwords")

	// 	c.JSON(200,
	// 		gin.H{
	// 			"data": d,
	// 		})
	// })

	r.POST("/IPs", func(c *gin.Context) {
		var data SubmitData

		// 绑定 JSON 请求体到 SubmitData 结构体
		if err := c.ShouldBindJSON(&data); err != nil {
			// 如果请求体有问题，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request data",
			})
			return
		}
		b, _ := utils.CheckPassword(data.Password)
		if !b {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Password Incorrect",
			})
			return
		}
		// 根据 action 执行不同的操作
		if data.Action == "overwrite" || data.Action == "append" {
			// 处理写文件
			SaveIps(data.Content, data.Action == "append")
		} else {
			// 如果 action 不符合要求，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid action",
			})
			return
		}

		// 返回成功的响应
		c.JSON(http.StatusOK, gin.H{
			"message": "Data submitted successfully",
		})
	})

	r.GET("/GetIPs", func(c *gin.Context) {
		Data := []byte(GetIPs())

		c.Data(http.StatusOK, "text/plain; charset=utf-8", Data)
	})

	r.POST("/EditPwd", func(c *gin.Context) {
		var data EditPassword

		// 绑定 JSON 请求体到 EditPassword 结构体
		if err := c.ShouldBindJSON(&data); err != nil {
			// 如果请求体有问题，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request data",
			})
			return
		}
		if data.NewPwd1 == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "新密码不能为空",
			})
			return
		}
		if data.NewPwd1 != data.NewPwd2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "两次输入密码不一致",
			})
			return
		}
		scuess, msg := utils.EditePassword(data.OldPwd, data.NewPwd1)
		if scuess {
			c.JSON(http.StatusOK, gin.H{
				"message": "两次输入密码不一致",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": msg,
			})
		}
	})
	r.POST("/TestHttpConnect", func(c *gin.Context) {
		var data TestHttpConnectData

		// 绑定 JSON 请求体到 EditPassword 结构体
		if err := c.ShouldBindJSON(&data); err != nil {
			// 如果请求体有问题，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request data",
			})
			return
		}
		ips, ok, msg := TestHttpConnect(data.IpText, data.TestDownload, data.Password)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"message": msg,
				"ips":     ips,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": msg,
				"ips":     ips,
			})
		}

	})

	r.GET("/GetConfig", func(c *gin.Context) {
		content, err := GetConfig()
		var Data []byte
		if err == nil {
			Data = []byte(content)

		} else {
			Data = []byte(err.Error())
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", Data)
	})

	r.POST("/SaveConfig", func(c *gin.Context) {
		var data ConfigData

		// 绑定 JSON 请求体到 EditPassword 结构体
		if err := c.ShouldBindJSON(&data); err != nil {
			// 如果请求体有问题，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request data",
			})
			return
		}

		b, _ := utils.CheckPassword(data.Password)
		if !b {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Password Incorrect",
			})
			return
		}
		ok := SaveConfig(data.Content)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"message": "保存成功，生效需要重启程序",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "保存失败",
			})
		}

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
		currentDownload, totalDownload, downloadIP, speed, downloadPort, downloadRemark, downloadDuration := GetProcessDownloadBar()
		currentDelay, totalDelay, delayIP, available, delayPort, delayRemark, delayDuration := GetProcessDelayBar()
		count := GetAllDataCount()
		ts := GetSchedules()

		message := gin.H{
			"Download": gin.H{
				"Current":  currentDownload,
				"Total":    totalDownload,
				"IP":       downloadIP,
				"Speed":    speed,
				"Port":     downloadPort,
				"Remark":   downloadRemark,
				"Duration": downloadDuration,
			},
			"Delay": gin.H{
				"Current":   currentDelay,
				"Total":     totalDelay,
				"IP":        delayIP,
				"Available": available,
				"Port":      delayPort,
				"Remark":    delayRemark,
				"Duration":  delayDuration,
			},
			"AllDataCount": count,
			"NextTime":     ts[0],
			"TraceInfo": gin.H{
				"Total": len(utils.Ips),
				"Index": utils.IpIndex,
			},
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
		if (currentDownload == totalDownload && currentDelay == totalDelay) || len(utils.Ips) != int(utils.IpIndex) {
			time.Sleep(1 * time.Second) // 推送间隔
		} else {
			time.Sleep(300 * time.Millisecond)
		}
	}
}
