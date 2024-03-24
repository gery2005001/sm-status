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

	// 获取当前状态的HTML字符串
	htmlData := "<h3>获取节点状态中......</h3>"
	// 向客户端发送数据
	if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
		log.Println("Failed to write message:", err)
	}
	//刷新节点状态
	htmlData = getPostStatusTableHTML()

	// 向客户端发送数据
	if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
		log.Println("WS Write failed:", err)
	}

	// 每隔指定时间推送状态
	ticker := time.NewTicker(config.Interval * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 获取状态数据
		htmlData := getPostStatusTableHTML()
		//log.Println(htmlData)
		// 向客户端发送数据
		if err := conn.WriteMessage(websocket.TextMessage, []byte(htmlData)); err != nil {
			log.Println("WS Write failed:", err)
			return
		}
		log.Println("WS Write successfully")
	}
}

// 根据config生成状态表HTML
func getPostStatusTableHTML() string {
	//输出Post状态表
	config := GetConfig()
	//config.refreshOperatorStatus()
	//log.Println(config.toJSONString())
	var htmlData string
	for n, node := range config.Node {
		if node.Enable {
			//获取node状态
			isSyncedText := ""
			stColor := ST_Failed_CSS
			if node.Status == ST_Success && node.IsSynced {
				isSyncedText = "已同步"
			} else {
				isSyncedText = "未同步"
			}
			stColor = config.Node[n].GetStatusColorCSS()
			//生成页面
			htmlData += "<table>"
			htmlData += "<colgroup><col class=\"st-column\"><col class=\"media-column\"><col class=\"small-column\"><col classe=\"auto-column\"><col classe=\"auto-column\"></colgroup>"
			htmlData += "<thead>"
			htmlData += "<tr><td class=\"td-left node-info\" colspan=\"5\">"
			htmlData += fmt.Sprintf("<span><b>状态：</b>"+"<span class=\"%s\">%s</span></span>", stColor, isSyncedText)
			htmlData += "<span><b>　Node名称：</b>" + node.Name + "</span>　<span><b>IP：</b>" + node.IP + "</span>　<span><b>版本：</b>" + config.Node[n].NodeVer + "</span>"
			htmlData += fmt.Sprintf("　<span><b><span>Peers：</b>%d</span>", config.Node[n].Peers)
			htmlData += fmt.Sprintf("　<span><b>Synced Layer：</b>%d</span>", config.Node[n].SLayer)
			htmlData += fmt.Sprintf("　<span><b>Top Layer：</b>%d</span>", config.Node[n].TLayer)
			htmlData += fmt.Sprintf("　<span><b>Verified Layer：</b>%d</span>", config.Node[n].VLayer)
			htmlData += "</td></tr>"
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
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	htmlData += "<b>更新时间:</b>" + currentTime + "</br>"
	htmlData += "<a href=\"/node\">切换到Node State</a></br>"

	return htmlData
}
