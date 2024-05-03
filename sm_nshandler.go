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
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade NS to WebSocket: ", err)
		return
	}
	defer conn.Close()

	// //刷新节点状态
	htmlData := getNodeStatusTableHTML()

	// 向客户端发送数据
	if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
		log.Println("NS WebSocket Write failed:", err)
	}
	log.Println("NS WebSocket Write successfully")

	// 每隔指定时间推送状态
	ticker := time.NewTicker(config.Interval * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if config.Updated {
			// 获取状态数据
			htmlData = getNodeStatusTableHTML()
			//log.Println(htmlData)
			// 向客户端发送数据
			if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
				log.Println("NS WebSocket Write failed:", err)
				return
			}
			log.Println("NS WebSocket Write successfully")
		}
	}
}

func getNodeStatusTableHTML() string {
	//输出Node 状态表
	config := GetConfig()

	htmlData := SmNetworkInfo.GetHtmlString()

	for n := 0; n < len(config.Node); n++ {
		htmlData += config.Node[n].GetNodeStatusTableHTMLString()
	}

	htmlData += fmt.Sprintf("latest version: <b>%s</b></br>", config.LatestVer)
	currentTime := config.UpdateTime.Format("2006-01-02 15:04:05")
	htmlData += "<b>更新时间:</b>" + currentTime + "</br>"
	htmlData += "<a href=\"/post\">切换到Post State</a></br>"
	htmlData += "<a href=\"/chunk\">切换到Chunks</a></br>"

	return htmlData
}
