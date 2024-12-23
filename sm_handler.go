package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sm-status/utility"
)

type PageData struct {
	RefreshTime   int
	StateContent  template.HTML
	FooterContent template.HTML
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("./static/index.html")
	if err != nil {
		log.Println("Open file error: ", err)
		http.Error(w, "Internal Server Error, HTML File Not Found", http.StatusInternalServerError)
	} else {
		http.ServeFile(w, r, "./static/index.html")
	}
}

func postStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 刷新Post状态
	refreshTime := int(appConfig.Refresh)
	footerHtml := getFooterHtml()
	htmlData := getPostStatusTableHTML()

	data := PageData{
		RefreshTime:   refreshTime,
		StateContent:  template.HTML(htmlData),
		FooterContent: template.HTML(footerHtml),
	}

	// 解析并执行 HTML 模板
	tmpl, err := template.ParseFiles("./static/template.gotmpl")
	if err != nil {
		http.Error(w, "无法解析模板文件", http.StatusInternalServerError)
		return
	}

	// 渲染模板并传递数据
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "渲染模板失败", http.StatusInternalServerError)
		return
	}
}

func nodeStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 刷新Node状态
	refreshTime := int(appConfig.Refresh)
	footerHtml := getFooterHtml()
	htmlData := getNodeStatusTableHTML()

	data := PageData{
		RefreshTime:   refreshTime,
		StateContent:  template.HTML(htmlData),
		FooterContent: template.HTML(footerHtml),
	}

	// 解析并执行 HTML 模板
	tmpl, err := template.ParseFiles("./static/template.gotmpl")
	if err != nil {
		http.Error(w, "无法解析模板文件", http.StatusInternalServerError)
		return
	}

	// 渲染模板并传递数据
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "渲染模板失败", http.StatusInternalServerError)
		return
	}

	log.Println("node status html update complete")
}

func chunkStatusHandler(w http.ResponseWriter, r *http.Request) {
	// //刷新Chunk状态
	refreshTime := int(appConfig.Refresh)
	footerHtml := getFooterHtml()
	htmlData := GetChunksTableHTML()

	data := PageData{
		RefreshTime:   refreshTime,
		StateContent:  template.HTML(htmlData),
		FooterContent: template.HTML(footerHtml),
	}

	// 解析并执行 HTML 模板
	tmpl, err := template.ParseFiles("./static/template.gotmpl")
	if err != nil {
		http.Error(w, "无法解析模板文件", http.StatusInternalServerError)
		return
	}

	// 渲染模板并传递数据
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "渲染模板失败", http.StatusInternalServerError)
		return
	}

	log.Println("chunks status html update complete")
}

func getFooterHtml() string {
	config := GetConfig()
	var reward float64 = 0
	htmlData := ""

	if RewardTotal > 0 {
		reward = float64(RewardTotal) / 1000000000
	}
	htmlData += "<div class=\"info-box\">"
	htmlData += fmt.Sprintf("<b>Total: </b> Units %d, Size  %s, Reward %.4f smh <br />", UnitTotal, utility.UnitsToTB(UnitTotal), reward)
	htmlData += fmt.Sprintf("<b>Latest version: </b>%s<br />", config.LatestVer)
	currentTime := config.UpdateTime.Format("2006-01-02 15:04:05")
	htmlData += "<b>Update Time: </b>" + currentTime + "<br /><br />"
	htmlData += "</div>"
	htmlData += "<a href=\"/node\"  class=\"link-button\">切换到Node State</a>"
	htmlData += "<a href=\"/post\"  class=\"link-button\">切换到Post State</a>"
	htmlData += "<a href=\"/chunk\"  class=\"link-button\">切换到Chunks</a><br />"

	return htmlData
}
