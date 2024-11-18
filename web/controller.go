package web

import (
	"github.com/XIU2/CloudflareSpeedTest/task"
)

func GetProcessDownloadBar() int64 {
	if task.DownloadBar == nil {
		return -1
	} else {
		current := task.DownloadBar.Current()
		return current
	}
}

func GetProcessDelayBar() int64 {
	if task.DownloadBar == nil {
		return -1
	} else {
		current := task.DelayBar.Current()
		return current
	}
}
