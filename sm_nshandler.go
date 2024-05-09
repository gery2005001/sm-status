package main

import (
	"fmt"
	"log"
	"net/http"
	"sm-status/utility"
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

	var reward float64 = 0
	if RewardTotal > 0 {
		reward = float64(RewardTotal) / 1000000000
	}
	htmlData += "<div class=\"info-box\">"
	htmlData += fmt.Sprintf("<b>Total: </b> Units %d, Size  %s, Reward %.4f smh <br />", UnitTotal, utility.UnitsToTB(UnitTotal), reward)
	htmlData += fmt.Sprintf("<b>Latest version: </b>%s<br />", config.LatestVer)
	currentTime := config.UpdateTime.Format("2006-01-02 15:04:05")
	htmlData += "<b>Update Time: </b>" + currentTime + "<br />"
	htmlData += "</div>"
	htmlData += "<a href=\"/post\"  class=\"link-button\">切换到Post State</a>"
	htmlData += "<a href=\"/chunk\"  class=\"link-button\">切换到Chunks</a><br />"

	return htmlData
}
