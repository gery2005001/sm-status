package main

import (
	"fmt"
)

// 根据config生成状态表HTML
func getPostStatusTableHTML() string {
	//输出Post状态表
	config := GetConfig()

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
			htmlData += "<colgroup><col class=\"st-column\"><col class=\"col-per-15\"><col class=\"col-per-10\"><col class=\"auto-column\"><col class=\"auto-column\"></colgroup>"
			htmlData += "<thead>"
			htmlData += "<tr><td class=\"node-info\" colspan=\"5\">"
			htmlData += "<span>状态：<b class=\"" + stColor + "\">" + isSyncedText + "</b></span>"
			htmlData += "<span>Node名称：<b>" + config.Node[n].Name + "</b></span>"
			htmlData += "<span>IP：<b>" + config.Node[n].IP + "</b></span>"
			htmlData += "<span>版本：<b class=\"" + verColor + "\">" + config.Node[n].NodeVer + "</b></span>"
			htmlData += fmt.Sprintf("<span>Peers：<b>%d</b></span>", config.Node[n].Peers)
			htmlData += fmt.Sprintf("<span>Synced Layer：<b>%d</b></span>", config.Node[n].SLayer)
			htmlData += fmt.Sprintf("<span>Top Layer：<b>%d</b></span>", config.Node[n].TLayer)
			htmlData += fmt.Sprintf("<span>Verified Layer：<b>%d</b></span>", config.Node[n].VLayer)
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
	return htmlData
}
