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
	//config.refreshNodeStatus()
	htmlData := ""
	for n := 0; n < len(config.Node); n++ {
		if config.Node[n].Enable {
			//获取node状态
			isSyncedText := ""
			stColor := config.Node[n].GetStatusColorCSS()
			if config.Node[n].Status == ST_Success && config.Node[n].IsSynced {
				isSyncedText = "已同步"
			} else {
				if config.Node[n].Status == ST_Empty {
					isSyncedText = "获取中"
					stColor = ST_Running_CSS
				} else {
					isSyncedText = "未同步"
				}
			}
			verColor := ""
			if config.Node[n].HasNewVer {
				verColor = ST_Failed_CSS
			}
			//生成页面
			htmlData += "<table>"
			htmlData += "<colgroup><col class=\"media-column\"><col class=\"auto-column\"><col class=\"media-column\"><col classe=\"auto-column\"><col classe=\"small-column\"></colgroup>"
			htmlData += "<thead>"
			htmlData += "<tr class=\"node-info\"><td class=\"td-left\" colspan=\"5\">"
			htmlData += fmt.Sprintf("<span><b>状态：</b>"+"<span class=\"%s\">%s</span></span>", stColor, isSyncedText)
			htmlData += "<span><b>　Node名称：</b>" + config.Node[n].Name + "</span>　<span><b>IP：</b>" + config.Node[n].IP + "</span><span><b>　版本：</b>"
			htmlData += fmt.Sprintf("<span class=\"%s\">%s</span></span>", verColor, config.Node[n].NodeVer)
			htmlData += fmt.Sprintf("　<span><b><span>Peers：</b>%d</span>", config.Node[n].Peers)
			htmlData += fmt.Sprintf("　<span><b>Synced Layer：</b>%d</span>", config.Node[n].SLayer)
			htmlData += fmt.Sprintf("　<span><b>Top Layer：%d</b></span>", config.Node[n].TLayer)
			htmlData += fmt.Sprintf("　<span><b>Verified Layer：</b>%d</span>", config.Node[n].VLayer)
			htmlData += "</td></tr>"
			htmlData += "<thead><tr><th>Name</th><th>ID</th><th>Eligibilities</th><th>State</th><th>Publish</th></tr></thead>"
			htmlData += "<tbody>"
			if config.Node[n].PostInfo != nil {
				for i := 0; i < len(config.Node[n].PostInfo); i++ {
					elgMsg := ""
					leftTime := ""
					bkColor := ""
					for _, elg := range config.Node[n].PostInfo[i].Eligs {
						if elg.Epoch >= config.Node[n].Epoch {
							if elg.Layer == config.Node[n].TLayer {
								leftTime = "now"
								bkColor = ST_Running_CSS
							} else if elg.Layer < config.Node[n].TLayer {
								lt := (config.Node[n].TLayer - elg.Layer) * SM_LayerDuration
								bkColor = ST_Success_CSS
								leftTime = "-" + utility.DurationToTimeFormat(time.Duration(lt)*time.Second)
							} else {
								lt := (elg.Layer - config.Node[n].TLayer) * SM_LayerDuration
								leftTime = utility.DurationToTimeFormat(time.Duration(lt) * time.Second)
							}
							elgMsg = fmt.Sprintf("<span class=\"%s\">【%s】</span>Layer:<b>%d</b>,Count:%d", bkColor, leftTime, elg.Layer, elg.Count)
						}
					}
					pwpMsg := ""
					if config.Node[n].PostInfo[i].Publish.Publish >= config.Node[n].Epoch {
						pwpMsg = fmt.Sprintf("Publish:%d,Target:%d", config.Node[n].PostInfo[i].Publish.Publish, config.Node[n].PostInfo[i].Publish.Target)
					}
					htmlData += fmt.Sprintf("<tr><td>%s</td><td class=\"td-rtl\">%x</td><td class=\"td-left\">%s</td><td class=\"td-left\">%s</td><td class=\"td-left\">%s</td><tr>", config.Node[n].PostInfo[i].Title, config.Node[n].PostInfo[i].SmesherId, elgMsg, config.Node[n].PostInfo[i].Status, pwpMsg)
				}
			}
			htmlData += "</tbody>"
			htmlData += "</table>"
		}
	}

	htmlData += fmt.Sprintf("latest version: <b>%s</b></br>", config.LatestVer)
	currentTime := config.UpdateTime.Format("2006-01-02 15:04:05")
	htmlData += "<b>更新时间:</b>" + currentTime + "</br>"
	htmlData += "<a href=\"/post\">切换到Post State</a></br>"

	//log.Println(htmlData)
	return htmlData
}
