package web

import (
	"time"

	"github.com/XIU2/CloudflareSpeedTest/task"
	"github.com/XIU2/CloudflareSpeedTest/utils"
)

var ()

func GetProcessDownloadBar() (current int64, total int64, ip interface{}, speed interface{}, port interface{}, remark interface{}, duration *utils.ProcessDuration) {
	duration = &utils.ProcessDuration{}
	duration.Hours = 0
	duration.Minutes = 0
	duration.Seconds = 0
	duration.StartTime = ""
	if task.DownloadBar == nil {
		return -1, -1, "", "", "", "", duration
	} else {
		current = task.DownloadBar.Current()
		total = task.DownloadBar.Total()
		ip = task.DownloadBar.GetOption("MyIP")
		speed = task.DownloadBar.GetOption("Speed")
		port = task.DownloadBar.GetOption("Port")
		remark = task.DownloadBar.GetOption("Remark")
		duration = task.DownloadBar.ProcessDuration()

		return current, total, ip, speed, port, remark, duration
	}
}

func GetProcessDelayBar() (current int64, total int64, ip interface{}, available interface{}, port interface{}, remark interface{}, duration *utils.ProcessDuration) {
	duration = &utils.ProcessDuration{}
	duration.Hours = 0
	duration.Minutes = 0
	duration.Seconds = 0
	duration.StartTime = ""
	if task.DelayBar == nil {
		return -1, -1, "", "", "", "", duration
	} else {
		current = task.DelayBar.Current()
		total = task.DelayBar.Total()
		ip = task.DelayBar.GetOption("MyIP")
		available = task.DelayBar.GetOption("MyStr")
		port = task.DelayBar.GetOption("Port")
		remark = task.DelayBar.GetOption("Remark")
		duration = task.DelayBar.ProcessDuration()
		return current, total, ip, available, port, remark, duration
	}
}

func GetAllDataCount() interface{} {
	sql := `select count(1) from speedTestResult`
	count, err := utils.Scalar(sql)
	if err != nil {
		return -1
	}
	return count

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
	select IP, Port,TraceInfo, MaxDownloadSpeed, MinDownloadSpeed, MinDelay, MaxDelay, AvgDelay, SumLossRate, AVGLossRate, Date, Count, Remark 
	from  MaxSpeed LIMIT 100
	`)
	if err != nil {
		// 处理错误
		if all == nil {

		}
	}
	return all
}

func GetYesterdayMaxData() []map[string]interface{} {
	tagert_day := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	sqlQuery := `
	SELECT IP, Port,max(TraceInfo)TraceInfo, MAX(MaxDownloadSpeed)MaxDownloadSpeed, MIN(MinDownloadSpeed)MinDownloadSpeed,ROUND(AVG(AvgDownloadSpeed),2)AvgDownloadSpeed, MIN(MinDelay)MinDelay, MAX(MaxDelay)MaxDelay, AVG(AvgDelay)AvgDelay, SUM(SumLossRate), AVG(AVGLossRate)AVGLossRate, Date, SUM(Count)Count, Remark 
	FROM MaxSpeed 
	WHERE Date = ? 
	GROUP BY IP, Port,Remark 
	order by SUM(Count)desc,AVG(AvgDownloadSpeed) desc,MIN(MinDownloadSpeed) desc,MAX(MaxDelay) asc,AVG(AVGLossRate) asc
	LIMIT 100`

	// 执行 SQL 查询，并传递 tagert_day 作为查询参数
	all, err := utils.Select(sqlQuery, tagert_day)
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
	SELECT IP, Port,max(TraceInfo)TraceInfo, MAX(MaxDownloadSpeed)MaxDownloadSpeed, MIN(MinDownloadSpeed)MinDownloadSpeed,ROUND(AVG(AvgDownloadSpeed),2)AvgDownloadSpeed, MIN(MinDelay)MinDelay, MAX(MaxDelay)MaxDelay, AVG(AvgDelay)AvgDelay, SUM(SumLossRate), AVG(AVGLossRate)AVGLossRate, Date, SUM(Count)Count, Remark 
	FROM MaxSpeed 
	WHERE Date = ? 
	GROUP BY IP, Port,Remark 
	order by SUM(Count)desc,AVG(AvgDownloadSpeed) desc,MIN(MinDownloadSpeed) desc,MAX(MaxDelay) asc,AVG(AVGLossRate) asc
	LIMIT 100`

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
	SELECT IP, Port,max(TraceInfo)TraceInfo, MAX(MaxDownloadSpeed)MaxDownloadSpeed, MIN(MinDownloadSpeed)MinDownloadSpeed,ROUND(AVG(AvgDownloadSpeed),2)AvgDownloadSpeed, MIN(MinDelay)MinDelay, MAX(MaxDelay)MaxDelay, AVG(AvgDelay)AvgDelay, SUM(SumLossRate), AVG(AVGLossRate)AVGLossRate, Date, SUM(Count)Count, Remark 
	FROM MaxSpeed 
	WHERE Date >= ? 
	GROUP BY IP, Port,Remark 
	order by SUM(Count)desc,AVG(AvgDownloadSpeed) desc,MIN(MinDownloadSpeed) desc,MAX(MaxDelay) asc,AVG(AVGLossRate) asc
	LIMIT 100`

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
	SELECT IP, Port,max(TraceInfo)TraceInfo, MAX(MaxDownloadSpeed)MaxDownloadSpeed, MIN(MinDownloadSpeed)MinDownloadSpeed,ROUND(AVG(AvgDownloadSpeed),2)AvgDownloadSpeed, MIN(MinDelay)MinDelay, MAX(MaxDelay)MaxDelay, AVG(AvgDelay)AvgDelay, SUM(SumLossRate), AVG(AVGLossRate)AVGLossRate, Date, SUM(Count)Count, Remark 
	FROM MaxSpeed 
	WHERE Date >= ? 
	GROUP BY IP, Port,Remark
	order by SUM(Count)desc,AVG(AvgDownloadSpeed) desc,MIN(MinDownloadSpeed) desc,MAX(MaxDelay) asc,AVG(AVGLossRate) asc
	LIMIT 100`

	// 执行 SQL 查询，并传递 tagert_day 作为查询参数
	all, err := utils.Select(sqlQuery, tagert_day)
	if err != nil {
		// 处理错误
		return nil
	}

	// 返回查询结果
	return all
}

func GetIPTraceInfos() []map[string]interface{} {
	selectSql := `select  IP,traceinfo from speedTestWithTrace where traceinfo is not null group by IP`
	// 执行 SQL 查询，并传递 tagert_day 作为查询参数
	all, err := utils.Select(selectSql)
	if err != nil {
		// 处理错误
		return nil
	}

	// 返回查询结果
	return all
}

func GetIPs() string {
	return utils.GetConfigFileContent()
}

// 保存IP信息
func SaveIps(contents string, append bool) {
	// IPFile
	utils.WriteToFile(utils.GetConfigIpFilePath(), contents, "utf-8", append)
}

func ClearTracerInfo() (b bool, e string) {
	_, err := utils.ExecNonQuery("delete from IpTraceInfos")
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func TraceInfosCount() (dataCount []map[string]interface{}, alldata []map[string]interface{}) {
	data, _ := utils.Select("select IP,count(*) Count from IpTraceInfos group by IP")
	data2, _ := utils.Select("select * from IpTraceInfos")
	return data, data2

}
