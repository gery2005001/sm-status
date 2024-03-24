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

var appConfig SmConfig = SmConfig{}
var configFile string = "config.json"

func GetTimeout() time.Duration {
	if appConfig.Timeout < time.Duration(3) || appConfig.Timeout > time.Duration(180) {
		return time.Duration(3)
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
		fmt.Println("Load config failed: ", err)
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&appConfig); err != nil {
		fmt.Println("Parse config failed: ", err)
		return err
	}

	appConfig.Ready = true

	log.Println("Load config successfully")

	return nil
}

// 获取当前APP配置
func GetConfig() *SmConfig {
	return &appConfig
}

// 刷新Post Operator
func (x *SmConfig) refreshOperatorStatus() {
	//刷新SmesherIDs
	for n := range x.Node {
		x.Node[n].GetPostOperatorStatus()
	}
}

// 刷新node status
func (x *SmConfig) refreshNodeStatus() {
	if x.Updated {
		currTime := time.Now()
		if currTime.Sub(x.UpdateTime) < time.Duration(300) {
			log.Println("skip network status update")
			return
		}
	}
	x.getLatestNodeVersion()
	//刷新SmesherIDs
	for n := range x.Node {
		x.Node[n].GetCurrentEpoch()
		x.Node[n].GetNodeStatus()
		x.Node[n].getPostInfoFromGRPC()
		x.Node[n].getEventsStreams()
	}
	x.Updated = true
	x.UpdateTime = time.Now()
}

// 从github获取node最新版本号
func (x *SmConfig) getLatestNodeVersion() {
	resp, err := http.Get(SM_GetNewVerAddress)
	if err != nil {
		log.Println("Ger new version failed: ", err)
	}
	defer resp.Body.Close()

	type Release struct {
		TagName string `json:"tag_name"`
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Println("Decode Json failed: ", err)
	}
	x.LatestVer = release.TagName

	log.Println("Successfully get latest version tag ", release.TagName)
}

// json格式输出config
// func (x *SmConfig) toJSONString() string {
// 	jsonData, err := json.MarshalIndent(x, " ", " ")
// 	if err != nil {
// 		return ""
// 	}

// 	return string(jsonData)
// }
