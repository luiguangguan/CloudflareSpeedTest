package utils

import (
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type Bar struct {
	pb        *pb.ProgressBar
	end       bool
	start     bool
	startTime time.Time
	endTime   time.Time
}

type ProcessDuration struct {
	Hours, Minutes, Seconds int
	StartTime               string
}

// var (
// 	end       bool = false
// 	start     bool = false
// 	startTime time.Time
// 	endTime   time.Time
// )

func NewBar(count int, MyStrStart, MyStrEnd string) *Bar {
	tmpl := fmt.Sprintf(`{{counters . }} {{ bar . "[" "-" (cycle . "↖" "↗" "↘" "↙" ) "_" "]"}} %s {{string . "MyStr" | green}} %s `, MyStrStart, MyStrEnd)
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	return &Bar{pb: bar, end: false, start: true, startTime: time.Now()}
}

func NewBar_httping(count int, MyStrStart, MyStrEnd string) *Bar {
	// 使用 %d 直接格式化状态码
	tmpl := fmt.Sprintf(`{{string . "MyIP" | yellow }}:{{string . "Port" | yellow }}{{string . "Option"}}：{{string . "OptionValue" | green}} {{counters . }} {{ bar . "[" "-" (cycle . "↖" "↗" "↘" "↙") "_" "]"}} %s {{string . "MyStr" | green}} %s`, MyStrStart, MyStrEnd)
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	return &Bar{pb: bar, end: false, start: true, startTime: time.Now()}
}

func NewBar_download(count int, MyStrStart, MyStrEnd string) *Bar {

	// 模板添加两行：第一行显示 IP 地址，第二行显示速度，第三行显示进度条
	tmpl := fmt.Sprintf(`{{string . "MyIP" | yellow }}:{{string . "Port" | yellow }}：{{string . "Speed" | green }}MB/s {{counters . }} {{ bar . "[" "-" (cycle . "↖" "↗" "↘" "↙") "_" "]"}} %s {{string . "MyStr" | green}} %s{{string . "Remark" | yellow }}`, MyStrStart, MyStrEnd)
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	return &Bar{pb: bar, end: false, start: true, startTime: time.Now()}
}

func (b *Bar) Grow(num int, MyStrVal string) {
	b.pb.Set("MyStr", MyStrVal).Add(num)
}

func (b *Bar) UpdateIPStatus(IP string, OptionValue int, port int, remark string) {
	strOptionValue := fmt.Sprintf("%d", OptionValue)
	strPort := fmt.Sprintf("%d", port)
	b.pb.Set("MyIP", IP).Set("OptionValue", strOptionValue).Set("Port", strPort).Set("Remark", remark)
}

func (b *Bar) UpdateIPSpeed(IP string, speed float64, port int, remark string) {
	strSpeed := fmt.Sprintf("%.2f", speed)
	strPort := fmt.Sprintf("%d", port)
	b.pb.Set("MyIP", IP).Set("Speed", strSpeed).Set("Port", strPort).Set("Remark", remark)
}

func (b *Bar) UpdateDownloadSpeed(speed float64) {
	strSpeed := fmt.Sprintf("%.2f", speed)
	b.pb.Set("Speed", strSpeed)
}

func (b *Bar) UpdateOption(Option string) {
	b.pb.Set("Option", Option)
}

func (b *Bar) Done() {
	b.pb.Finish()
	b.end = true
	b.endTime = time.Now()
}

func (b *Bar) Current() int64 {
	return b.pb.Current()
}

func (b *Bar) Total() int64 {
	return b.pb.Total()
}
func (b *Bar) GetOption(key string) interface{} {
	if b == nil {
		return ""
	} else {
		return b.pb.Get(key)
	}
}

func (b *Bar) ProcessDuration() (Pduration *ProcessDuration) {
	Pduration = &ProcessDuration{}
	Pduration.StartTime = ""
	if b.start {
		Pduration.StartTime = b.startTime.Format("2006-01-02 15:04:05")
		var duration time.Duration
		if b.end {
			duration = b.endTime.Sub(b.startTime)
		} else {
			duration = time.Since(b.startTime)
		}

		Pduration.Hours = int(duration.Hours())
		Pduration.Minutes = int(duration.Minutes()) % 60
		Pduration.Seconds = int(duration.Seconds()) % 60
		return Pduration
	}
	return Pduration
}
