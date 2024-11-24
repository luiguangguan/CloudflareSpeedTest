package utils

import (
	"time"

	"github.com/robfig/cron/v3"
)

var (
	CronExpr string
)

func GetScheduleTimes(expr string, count int) ([]time.Time, error) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	// 解析 Cron 表达式
	schedule, err := parser.Parse(expr)
	if err != nil {
		return nil, err
	}

	// 获取当前时间
	now := time.Now()
	var times []time.Time

	// 迭代计算未来的时间点
	for i := 0; i < count; i++ {
		next := schedule.Next(now)
		times = append(times, next)
		now = next
	}

	return times, nil
}
