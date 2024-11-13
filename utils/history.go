package utils

import (
	"math"
	"time"
)

func Save(data *CloudflareIPData) int64 {
	currentTime := time.Now()

	// 截斷保留兩位小數
	downloadSpeed := math.Trunc((data.DownloadSpeed/1024/1024)*100) / 100 // 直接截断至两位小数
	// fmt.Println(value)                  // 输出: 123.45

	// 格式化时间为 "yyyy-MM-dd HH:mm:ss" 格式
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	insert := `
	INSERT INTO speedTestResult (IP, Port, Sended, Received, Delay, LossRate, DownloadSpeed, CreateTime,Remark) VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?);`
	rows, _ := ExecNonQuery(insert, data.IP.String(), data.Port, data.Sended, data.Received, data.Delay.Seconds()*1000, data.LossRate(), downloadSpeed, formattedTime, data.Remark)
	return rows
}
