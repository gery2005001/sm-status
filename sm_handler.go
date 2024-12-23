package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageData struct {
	RefreshTime  int
	StateContent template.HTML
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
	// //刷新Post状态
	htmlData := getPostStatusTableHTML()
	refreshTime := int(appConfig.Refresh)

	data := PageData{
		RefreshTime:  refreshTime,
		StateContent: template.HTML(htmlData),
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
	// //刷新Node状态
	htmlData := getNodeStatusTableHTML()
	refreshTime := int(appConfig.Refresh)

	data := PageData{
		RefreshTime:  refreshTime,
		StateContent: template.HTML(htmlData),
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

func chunkStatusHandler(w http.ResponseWriter, r *http.Request) {
	// //刷新Chunk状态
	htmlData := getAllChunksTableHTML()
	refreshTime := int(appConfig.Refresh)

	data := PageData{
		RefreshTime:  refreshTime,
		StateContent: template.HTML(htmlData),
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
