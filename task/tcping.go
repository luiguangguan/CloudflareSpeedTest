package task

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/utils"
)

const (
	tcpConnectTimeout = time.Second * 1
	maxRoutine        = 1000
	defaultRoutines   = 200
	defaultPort       = 443
	defaultPingTimes  = 4
)

var (
	Routines      = defaultRoutines
	TCPPort   int = defaultPort
	PingTimes int = defaultPingTimes
	DelayBar  *utils.Bar
)

type Ping struct {
	wg      *sync.WaitGroup
	m       *sync.Mutex
	ips     []*net.IPAddr
	ips2    []*IPPort
	csv     utils.PingDelaySet
	control chan bool
	bar     *utils.Bar
}

func checkPingDefault() {
	if Routines <= 0 {
		Routines = defaultRoutines
	}
	if TCPPort <= 0 || TCPPort >= 65535 {
		TCPPort = defaultPort
	}
	if PingTimes <= 0 {
		PingTimes = defaultPingTimes
	}
}

func NewPing() *Ping {
	checkPingDefault()
	ips2 := loadIPRanges()
	DelayBar = utils.NewBar_httping(len(ips2), "可用:", "")
	return &Ping{
		wg: &sync.WaitGroup{},
		m:  &sync.Mutex{},
		// ips:     ips,
		ips2:    ips2,
		csv:     make(utils.PingDelaySet, 0),
		control: make(chan bool, Routines),
		bar:     DelayBar,
	}
}

func (p *Ping) Run() utils.PingDelaySet {
	if len(p.ips2) == 0 {
		return p.csv
	}
	if Httping {
		fmt.Printf("开始延迟测速（模式：HTTP, 端口：%d, 范围：%v ~ %v ms, 丢包：%.2f)\n", TCPPort, utils.InputMinDelay.Milliseconds(), utils.InputMaxDelay.Milliseconds(), utils.InputMaxLossRate)
		p.bar.UpdateOption("状态码")
	} else {
		fmt.Printf("开始延迟测速（模式：TCP, 端口：%d, 范围：%v ~ %v ms, 丢包：%.2f)\n", TCPPort, utils.InputMinDelay.Milliseconds(), utils.InputMaxDelay.Milliseconds(), utils.InputMaxLossRate)
		p.bar.UpdateOption("延迟")
	}
	for _, ip := range p.ips2 {
		p.wg.Add(1)
		p.control <- false
		go p.start(ip.IP, ip.Port, ip.Remark)
	}
	p.wg.Wait()
	p.bar.Done()
	sort.Sort(p.csv)
	return p.csv
}

func (p *Ping) start(ip *net.IPAddr, port int, remark string) {
	defer p.wg.Done()
	p.tcpingHandler(ip, port, remark)
	<-p.control
}

// bool connectionSucceed float32 time
func (p *Ping) tcping(ip *net.IPAddr, port int) (bool, time.Duration) {
	startTime := time.Now()
	var fullAddress string
	var targetPort int = TCPPort
	if port > 0 {
		targetPort = port
	}

	if isIPv4(ip.String()) {
		fullAddress = fmt.Sprintf("%s:%d", ip.String(), targetPort)
	} else {
		fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), targetPort)
	}
	conn, err := net.DialTimeout("tcp", fullAddress, tcpConnectTimeout)
	if err != nil {
		return false, 0
	}
	defer conn.Close()
	duration := time.Since(startTime)
	return true, duration
}

// pingReceived pingTotalTime
func (p *Ping) checkConnection(ip *net.IPAddr, port int) (recv int, totalDelay time.Duration) {
	if Httping {
		recv, totalDelay = p.httping(ip, port)
		return
	}
	for i := 0; i < PingTimes; i++ {
		if ok, delay := p.tcping(ip, port); ok {
			recv++
			totalDelay += delay
			//实时延迟
			p.bar.UpdateIPStatus(ip.String(), int(delay.Milliseconds()))
		}
	}
	return
}

func (p *Ping) appendIPData(data *utils.PingData) {
	p.m.Lock()
	defer p.m.Unlock()
	p.csv = append(p.csv, utils.CloudflareIPData{
		PingData: data,
	})
}

// handle tcping
func (p *Ping) tcpingHandler(ip *net.IPAddr, port int, remark string) {
	recv, totalDlay := p.checkConnection(ip, port)
	nowAble := len(p.csv)
	if recv != 0 {
		nowAble++
	}
	p.bar.Grow(1, strconv.Itoa(nowAble))
	if recv == 0 {
		return
	}
	data := &utils.PingData{
		IP:       ip,
		Sended:   PingTimes,
		Received: recv,
		Delay:    totalDlay / time.Duration(recv),
		Port:     port,
		Remark:   remark,
	}
	p.appendIPData(data)
}
