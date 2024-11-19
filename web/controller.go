package web

import (
	"github.com/XIU2/CloudflareSpeedTest/task"
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
