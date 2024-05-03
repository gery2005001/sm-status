package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Post Status页面处理
func postStatusWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	config := GetConfig()

	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade PS to WebSocket: ", err)
		return
	}
	defer conn.Close()

	// // 获取当前状态的HTML字符串
	htmlData := getPostStatusTableHTML()
	// 向客户端发送数据
	if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
		log.Println("PS WebSocket Write failed:", err)
	}
	log.Println("PS WebSocket Write successfully")

	// 每隔指定时间推送状态
	ticker := time.NewTicker(config.Interval * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		htmlData = getPostStatusTableHTML()
		//log.Println(htmlData)
		// 向客户端发送数据
		if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
			log.Println("PS WebSocket Write failed:", err)
			return
		}
		log.Println("PS WebSocket Write successfully")
	}
}

// 根据config生成状态表HTML
func getPostStatusTableHTML() string {
	//输出Post状态表
	config := GetConfig()
	//config.refreshOperatorStatus()
	//log.Println(config.toJSONString())

	htmlData := SmNetworkInfo.GetHtmlString()

	for n, node := range config.Node {
		if node.Enable {
			//获取node状态
			isSyncedText := ""
			stColor := config.Node[n].GetStatusColorCSS()
			if config.Node[n].Status == ST_Success && config.Node[n].IsSynced {
				isSyncedText = "【已同步】"
			} else {
				if config.Node[n].Status == ST_Empty {
					isSyncedText = "【获取中】"
					stColor = ST_Running_CSS
				} else {
					isSyncedText = "【未同步】"
				}
			}
			verColor := ""
			if config.Node[n].HasNewVer {
				verColor = ST_Failed_CSS
			}
			//生成页面
			htmlData += "<table>"
			htmlData += "<colgroup><col class=\"st-column\"><col class=\"col-per-15\"><col class=\"col-per-10\"><col classe=\"auto-column\"><col classe=\"auto-column\"></colgroup>"
			htmlData += "<thead>"
			htmlData += "<tr><td class=\"td-left node-info\" colspan=\"5\">"
			htmlData += fmt.Sprintf("<span>状态：<b>"+"<span class=\"%s\">%s</span></b></span>", stColor, isSyncedText)
			htmlData += "<span>　Node名称：<b>" + config.Node[n].Name + "</b></span>　<span>IP：<b>" + config.Node[n].IP + "</b></span>"
			htmlData += fmt.Sprintf("<span>　版本：<span class=\"%s\"><b>%s</b></span></span>", verColor, config.Node[n].NodeVer)
			htmlData += fmt.Sprintf("　<span><span>Peers：<b>%d</b></span>", config.Node[n].Peers)
			htmlData += fmt.Sprintf("　<span>Synced Layer：<b>%d</b></span>", config.Node[n].SLayer)
			htmlData += fmt.Sprintf("　<span>Top Layer：<b>%d</b></span>", config.Node[n].TLayer)
			htmlData += fmt.Sprintf("　<span>Verified Layer：<b>%d</b></span>", config.Node[n].VLayer)
			htmlData += "</td></tr>"
			htmlData += "</thead>"
			htmlData += "<thead><tr><th>ST</th><th>名称</th><th>容量</th><th>Operator</th><th>OperatorAddress</th></tr></thead>"

			if len(node.Post) > 0 {
				htmlData += "<tbody>"
				for _, post := range node.Post {
					stColor := post.GetStatusColorCSS()
					htmlData += fmt.Sprintf("<tr><td class=\"%s\">%s</td><td>%s</td><td>%s</td><td class=\"td-left\">%s</td><td class=\"td-left\">%s</td></tr>", stColor, post.Status, post.Title, post.Capacity, post.OaStatus, post.OperatorAddress)
				}
				htmlData += "</tbody>"
			}
			htmlData += "</table>"
		}
	}
	htmlData += fmt.Sprintf("latest version: <b>%s</b></br>", config.LatestVer)
	currentTime := config.UpdateTime.Format("2006-01-02 15:04:05")
	htmlData += "<b>更新时间:</b>" + currentTime + "</br>"
	htmlData += "<a href=\"/node\">切换到Node State</a></br>"

	return htmlData
}
