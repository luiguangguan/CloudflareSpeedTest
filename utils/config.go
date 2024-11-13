package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	Routines          int     `json:"n"`
	PingTimes         int     `json:"t"`
	TestCount         int     `json:"dn"`
	DownloadTime      int     `json:"dt"`
	TCPPort           int     `json:"tp"`
	URL               string  `json:"url"`
	Httping           bool    `json:"httping"`
	HttpingStatusCode int     `json:"httping-code"`
	HttpingCFColo     string  `json:"cfcolo"`
	MaxDelay          int     `json:"tl"`
	MinDelay          int     `json:"tll"`
	MaxLossRate       float64 `json:"tlr"`
	MinSpeed          float64 `json:"sl"`
	PrintNum          int     `json:"p"`
	IPFile            string  `json:"f"`
	IPText            string  `json:"ip"`
	Output            string  `json:"o"`
	Disable           bool    `json:"dd"`
	TestAll           bool    `json:"allip"`
	DbFile            string  `json:"db"`
	CronExpr          string  `json:cron`
}

// LoadConfig reads a JSON configuration file and returns a Config struct
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	// 使用文件中的配置值更新 flag 值
	return &config, nil
}
