package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sm-status/version"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func init() {
	fmt.Println(version.PrintCLIVersion())
}

func main() {
	//获取命令行参数
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if (args[i] == "--config" || args[i] == "-c") && i+1 < len(args) {
			configFile = args[i+1]
			break
		}
	}

	//加载配置文件
	if !appConfig.Ready {
		err := LoadConfig()
		if err != nil {
			log.Fatal("Application Exit")
		}
	}

	go appConfig.refreshNodeStatus()

	// 定时刷新Node和Post状态
	ticker := time.NewTicker(appConfig.Refresh * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			appConfig.refreshNodeStatus()
		}
	}()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/post", postStatusHandler)
	http.HandleFunc("/node", nodeStatusHandler)

	http.HandleFunc("/ps", postStatusWebSocketHandler)
	http.HandleFunc("/ns", nodeStatusWebSocketHandler)

	log.Println("Server started at port", appConfig.Port)
	port := fmt.Sprintf(":%d", appConfig.Port)
	log.Fatal(http.ListenAndServe(port, nil))

}
