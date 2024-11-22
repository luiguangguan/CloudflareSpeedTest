package web

import (
	"github.com/XIU2/CloudflareSpeedTest/task"
	"github.com/XIU2/CloudflareSpeedTest/utils"
	"time"
)

var (
	ScheduleTime []time.Time
)

func GetProcessDownloadBar() (current int64, total int64) {
	if task.DownloadBar == nil {
		return -1, -1
	} else {
		current := task.DownloadBar.Current()
		total := task.DownloadBar.Total()
		return current, total
	}
}

func GetProcessDelayBar() (current int64, total int64) {
	if task.DelayBar == nil {
		return -1, -1
	} else {
		current := task.DelayBar.Current()
		total := task.DelayBar.Total()
		return current, total
	}
}

func GetSchedules() []string {
	var ts []string
	for _, t := range ScheduleTime {
		localTime := t.Local() // 转为本地时区时间
		ts = append(ts, localTime.Format("2006-01-02 15:04:05"))
	}
	return ts
}

func GetAllData() []map[string]interface{} {
	a, err := utils.Select("select * from speedTestResult")
	if err != nil {
		// 处理错误
		if a == nil {

		}
	}

	return a
}
