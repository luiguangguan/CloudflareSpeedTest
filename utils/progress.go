package utils

import (
	"fmt"

	"github.com/cheggaaa/pb/v3"
)

type Bar struct {
	pb *pb.ProgressBar
}

func NewBar(count int, MyStrStart, MyStrEnd string) *Bar {
	tmpl := fmt.Sprintf(`{{counters . }} {{ bar . "[" "-" (cycle . "↖" "↗" "↘" "↙" ) "_" "]"}} %s {{string . "MyStr" | green}} %s `, MyStrStart, MyStrEnd)
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	return &Bar{pb: bar}
}

func NewBar_httping(count int, MyStrStart, MyStrEnd string) *Bar {
	// 使用 %d 直接格式化状态码
	tmpl := fmt.Sprintf(`{{string . "MyIP" | yellow }}  {{string . "Option"}}：{{string . "OptionValue" | green}} {{counters . }} {{ bar . "[" "-" (cycle . "↖" "↗" "↘" "↙") "_" "]"}} %s {{string . "MyStr" | green}} %s`, MyStrStart, MyStrEnd)
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	return &Bar{pb: bar}
}

func NewBar_download(count int, MyStrStart, MyStrEnd string) *Bar {

	// 模板添加两行：第一行显示 IP 地址，第二行显示速度，第三行显示进度条
	tmpl := fmt.Sprintf(`{{string . "MyIP" | yellow }}：{{string . "Speed" | green }}MB/s {{counters . }} {{ bar . "[" "-" (cycle . "↖" "↗" "↘" "↙") "_" "]"}} %s {{string . "MyStr" | green}} %s`, MyStrStart, MyStrEnd)
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	return &Bar{pb: bar}
}

func (b *Bar) Grow(num int, MyStrVal string) {
	b.pb.Set("MyStr", MyStrVal).Add(num)
}

func (b *Bar) UpdateIPStatus(IP string, OptionValue int) {
	strOptionValue := fmt.Sprintf("%d", OptionValue)
	b.pb.Set("MyIP", IP).Set("OptionValue", strOptionValue)
}

func (b *Bar) UpdateIPSpeed(IP string, speed float64) {
	strSpeed := fmt.Sprintf("%.2f", speed)
	b.pb.Set("MyIP", IP).Set("Speed", strSpeed)
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
