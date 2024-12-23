package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sm-status/version"
	"time"

	"github.com/urfave/cli/v2"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "config",
		Usage:    "Config file path in format <dir>/<filename>.json",
		Required: false,
		Value:    "config.json",
		EnvVars:  []string{"SM_STATUS_CONFIG_PATH"},
	},
}

func init() {
	fmt.Println(version.PrintCLIVersion())
}

func main() {
	app := cli.NewApp()
	app.Name = "sm-status"
	app.Version = version.PrintCLIVersionNumber()
	app.Flags = flags

	app.Action = func(ctx *cli.Context) error {
		if ctx.IsSet("config") {
			configFile = ctx.String("config")
		}
		if err := LoadConfig(); err != nil {
			return fmt.Errorf("%w", err)
		}

		//刷新Node和Post状态
		go appConfig.refreshNodeStatus()

		// 设置定时任务刷新Node和Post状态
		ticker := time.NewTicker(appConfig.Refresh * time.Second)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				log.Println("Tick received")
				//更新Config文件
				if appConfig.Reload {
					log.Println("Reloading config...")
					_ = LoadConfig()
				}
				log.Println("Refreshing status...")
				appConfig.refreshNodeStatus()
			}
		}()

		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		http.HandleFunc("/", rootHandler)
		http.HandleFunc("/post", postStatusHandler)
		http.HandleFunc("/node", nodeStatusHandler)
		http.HandleFunc("/chunk", chunkStatusHandler)

		port := fmt.Sprintf(":%d", appConfig.Port)
		log.Printf("Server started at port %d", appConfig.Port)

		if err := http.ListenAndServe(port, nil); err != nil {
			return fmt.Errorf("%w", err)
		}

		log.Println("server is shutdown")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}

}
