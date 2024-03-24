package main

import (
	"encoding/json"
	"fmt"
	"log"
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

// json格式输出config
// func (x *SmConfig) toJSONString() string {
// 	jsonData, err := json.MarshalIndent(x, " ", " ")
// 	if err != nil {
// 		return ""
// 	}

// 	return string(jsonData)
// }
