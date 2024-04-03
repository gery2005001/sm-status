package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type SmConfig struct {
	Port       int           `json:"port"`
	Refresh    time.Duration `json:"refresh"`
	Interval   time.Duration `json:"interval"`
	Timeout    time.Duration `json:"timeout"`
	Reload     bool          `json:"reload"`
	Node       []Node        `json:"node"`
	LatestVer  string
	Ready      bool
	Updated    bool
	UpdateTime time.Time
}

const MIN_TIMEOUT = 3
const MAX_TIMEOUT = 30
const MIN_REFRESH_TIME = 15
const MAX_REFRESH_TIME = 600

var appConfig SmConfig = SmConfig{}
var configFile string = "config.json"

func GetTimeout() time.Duration {
	if appConfig.Timeout < time.Duration(MIN_TIMEOUT) || appConfig.Timeout > time.Duration(MAX_TIMEOUT) {
		return time.Duration(MIN_TIMEOUT)
	}
	return appConfig.Timeout
}

func GetLatestVer() string {
	return appConfig.LatestVer
}

// 加载配置文件
func LoadConfig() error {
	// 打开 config 文件
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Println("load config file failed.")
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&appConfig); err != nil {
		fmt.Println("parse config failed.")
		return err
	}

	if appConfig.Refresh < time.Duration(MIN_REFRESH_TIME) || appConfig.Timeout > time.Duration(MAX_REFRESH_TIME) {
		appConfig.Refresh = time.Duration(MIN_REFRESH_TIME)
	}

	appConfig.Ready = true

	log.Println("load config successfully")
	return nil
}

// 获取当前APP配置
func GetConfig() *SmConfig {
	return &appConfig
}

// 刷新node status
func (x *SmConfig) refreshNodeStatus() {
	if x.Updated {
		currTime := time.Now()
		if currTime.Sub(x.UpdateTime) < time.Duration(SM_LayerDuration) {
			log.Println("skip status update...")
			return
		}
		log.Println(currTime.Sub(x.UpdateTime), "have passed")
	}

	if appConfig.Reload {
		LoadConfig()
	}

	//获取最新的客户端版本
	x.getLatestNodeVersion()
	//刷新每个Node的Post和Operator状态
	for n := range x.Node {
		x.Node[n].GetAllNodeInformation()
		x.Node[n].getNodePostOperatorStatus()
	}

	x.Updated = true
	x.UpdateTime = time.Now()
}

// 从github获取node最新版本号
func (x *SmConfig) getLatestNodeVersion() {
	resp, err := http.Get(SM_GetNewVerAddress)
	if err != nil {
		log.Println("get new version failed: ", err)
	}
	defer resp.Body.Close()

	type Release struct {
		TagName string `json:"tag_name"`
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Println("decode Json failed: ", err)
	}
	x.LatestVer = release.TagName

	log.Println("successfully get latest version tag ", release.TagName)
}
