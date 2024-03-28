package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Node Status页面处理
func nodeStatusWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	config := GetConfig()
	if config.Reload {
		err := LoadConfig()
		if err != nil {
			log.Println("Reload config error: ", err)
		}
	}
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket: ", err)
		return
	}
	defer conn.Close()

	// // 获取当前状态的HTML字符串
	// htmlData := "<h3>获取节点状态中......</h3>"
	// // 向客户端发送数据
	// if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
	// 	log.Println("Failed to write message:", err)
	// }
	// //刷新节点状态
	htmlData := getNodeStatusTableHTML()

	// 向客户端发送数据
	if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
		log.Println("WS Write failed:", err)
	}

	// 每隔指定时间推送状态
	ticker := time.NewTicker(config.Interval * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 获取状态数据
		htmlData = getNodeStatusTableHTML()
		//log.Println(htmlData)
		// 向客户端发送数据
		if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
			log.Println("WS Write failed:", err)
			return
		}
		log.Println("WS Write successfully")
	}
}

func getNodeStatusTableHTML() string {
	//输出Node 状态表
	config := GetConfig()

	// for _, nd := range config.Node {
	// 	for _, pi := range nd.PostInfo {
	// 		jd, err := json.MarshalIndent(pi, " ", " ")
	// 		if err == nil {
	// 			log.Printf("%s", string(jd))
	// 		} else {
	// 			log.Println("marshal post info to json error ", err)
	// 		}
	// 	}
	// }

	htmlData := ""
	for n := 0; n < len(config.Node); n++ {
		htmlData += config.Node[n].GetNodeStatusTableHTMLString()
	}

	htmlData += fmt.Sprintf("latest version: <b>%s</b></br>", config.LatestVer)
	currentTime := config.UpdateTime.Format("2006-01-02 15:04:05")
	htmlData += "<b>更新时间:</b>" + currentTime + "</br>"
	htmlData += "<a href=\"/post\">切换到Post State</a></br>"

	return htmlData
}
