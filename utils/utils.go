package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	CronExpr         string
	runningFunctions int32
	TraceRunning     = false
	Ips              []string
	IpIndex          int32
	saveTraceMutex   sync.Mutex
)

const MaxConcurrent = 3 // 最大同时路由跟踪任务数

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

// TraceRoute 函数，执行路由跟踪
func TraceRoute(ip string) (string, error) {
	var cmd *exec.Cmd

	// 根据操作系统选择命令
	if runtime.GOOS == "windows" {
		// Windows 使用 tracert，并设置 chcp 65001 强制为 UTF-8 编码
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("chcp 65001 >nul & tracert %s", ip))
	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		cmd = exec.Command("traceroute", ip) // Linux 和 macOS 使用 traceroute
	} else {
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute %s: %v\nOutput: %s", cmd.Path, err, string(output))
	}

	// 移除 Windows 命令中可能的额外空行
	cleanOutput := strings.ReplaceAll(string(output), "\r\n", "\n")

	return cleanOutput, nil
}

// func TraceRoute2(ip string) (string, error) {
// 	// Configure traceroute options
// 	options := traceroute.TracerouteOptions{
// 		MaxHops:    30,  // Maximum hops
// 		TimeoutMs:  500, // Timeout in milliseconds
// 		PacketSize: 52,  // Packet size in bytes
// 	}

// 	// Perform the traceroute
// 	results, err := traceroute.Traceroute(ip, &options)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to perform traceroute to %s: %w", ip, err)
// 	}

// 	// Format the results
// 	var builder strings.Builder
// 	for _, hop := range results {
// 		builder.WriteString(fmt.Sprintf("%2d  %s  %v\n", hop.TTL, hop.Address, hop.ElapsedTime))
// 	}

// 	return builder.String(), nil
// }

func TraceIP(ip string) {
	atomic.AddInt32(&runningFunctions, 1) // 增加计数器
	defer func() {
		atomic.AddInt32(&runningFunctions, -1) // 确保在函数结束时减少计数器
		atomic.AddInt32(&IpIndex, 1)           // 添加已經處理IP的計數
	}()
	result, err := TraceRoute(ip) // 调用 TraceRoute 函数执行路由跟踪
	if err == nil {
		if result != "" {
			// 加锁，确保只有一个线程在执行 SaveTrace
			saveTraceMutex.Lock()
			row := SaveTrace(ip, result, "GZ")
			saveTraceMutex.Unlock()
			if row > 0 {
				// fmt.Printf("保存成功\n")
			}
		}
	}
}

func TraceRouteIP() {
	defer func() {
		TraceRunning = false
	}()
	TraceRunning = true
	for {
		// whiile循環
		Ips := GetAllIPNoneTrace()
		IpIndex = 0
		if Ips != nil {
			for _, ip := range Ips {
				// 检查计数器，如果已经有太多任务在运行，则等待
				for atomic.LoadInt32(&runningFunctions) >= MaxConcurrent {
					time.Sleep(2 * time.Second) // 等待 1 秒钟
				}
				go TraceIP(ip) // 启动新的 goroutine 执行路由跟踪
			}
		}
	}
}
