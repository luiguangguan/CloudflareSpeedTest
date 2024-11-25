package web

import (
	"time"

	"github.com/XIU2/CloudflareSpeedTest/task"
	"github.com/XIU2/CloudflareSpeedTest/utils"
)

var ()

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
	if utils.CronExpr != "" {
		scheduleTime, err := utils.GetScheduleTimes(utils.CronExpr, 20)
		if err != nil {
			ts = append(ts, "Error")
			return ts
		}
		for _, t := range scheduleTime {
			localTime := t.Local() // 转为本地时区时间
			ts = append(ts, localTime.Format("2006-01-02 15:04:05"))
		}
		return ts
	} else {
		return ts
	}

}

func GetAllData() []map[string]interface{} {
	all, err := utils.Select("select * from speedTestResult")
	if err != nil {
		// 处理错误
		if all == nil {

		}
	}
	return all
}

func GetMaxData() []map[string]interface{} {
	all, err := utils.Select(`
	select IP, Port, MaxDownloadSpeed, MinDownloadSpeed, MinDelay, MaxDelay, AvgDelay, SumLossRate, AVGLossRate, Date, Count, Remark 
	from  MaxSpeed LIMIT 100
	`)
	if err != nil {
		// 处理错误
		if all == nil {

		}
	}
	return all
}

func Get1DayMaxData() []map[string]interface{} {
	tagert_day := time.Now().Format("2006-01-02")
	sqlQuery := `
	SELECT IP, Port, MaxDownloadSpeed, MinDownloadSpeed, MinDelay, MaxDelay, AvgDelay, SumLossRate, AVGLossRate, Date, Count, Remark 
	FROM MaxSpeed 
	WHERE Date = ?  LIMIT 50`

	// 执行 SQL 查询，并传递 tagert_day 作为查询参数
	all, err := utils.Select(sqlQuery, tagert_day)
	if err != nil {
		// 处理错误
		if all == nil {

		}
	}
	return all
}

func Get3DayMaxData() []map[string]interface{} {
	// 获取当前日期（忽略时间部分）
	tagert_day := time.Now().AddDate(0, 0, -3).Format("2006-01-02") // 获取三天前的日期，并格式化为 "YYYY-MM-DD"

	// 使用 SQL 查询过去3天的数据，条件使用字符串格式的日期进行比较
	sqlQuery := `
		SELECT IP, Port, MaxDownloadSpeed, MinDownloadSpeed, MinDelay, MaxDelay, AvgDelay, SumLossRate, AVGLossRate, Date, Count, Remark 
		FROM MaxSpeed 
		WHERE Date >= ? LIMIT 50`

	// 执行 SQL 查询，并传递 tagert_day 作为查询参数
	all, err := utils.Select(sqlQuery, tagert_day)
	if err != nil {
		// 处理错误
	}

	// 返回查询结果
	return all
}

func Get5DayMaxData() []map[string]interface{} {
	// 获取当前日期（忽略时间部分）
	tagert_day := time.Now().AddDate(0, 0, -5).Format("2006-01-02") // 获取五天前的日期，并格式化为 "YYYY-MM-DD"

	// 使用 SQL 查询过去3天的数据，条件使用字符串格式的日期进行比较
	sqlQuery := `
		SELECT IP, Port, MaxDownloadSpeed, MinDownloadSpeed, MinDelay, MaxDelay, AvgDelay, SumLossRate, AVGLossRate, Date, Count, Remark 
		FROM MaxSpeed 
		WHERE Date >= ? LIMIT 50`

	// 执行 SQL 查询，并传递 tagert_day 作为查询参数
	all, err := utils.Select(sqlQuery, tagert_day)
	if err != nil {
		// 处理错误
	}

	// 返回查询结果
	return all
}
