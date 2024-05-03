package main

import (
	"log"
	"net/http"
	"os"
)

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
	_, err := os.Stat("./static/poststatus.html")
	if err != nil {
		log.Println("Error serving HTML file:", err)
		http.Error(w, "Internal Server Error, HTML File Not Found", http.StatusInternalServerError)
	} else {
		http.ServeFile(w, r, "./static/poststatus.html")
	}
}

func nodeStatusHandler(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("./static/nodestatus.html")
	if err != nil {
		log.Println("Error serving HTML file:", err)
		http.Error(w, "Internal Server Error, HTML File Not Found", http.StatusInternalServerError)
	} else {
		http.ServeFile(w, r, "./static/nodestatus.html")
	}
}

func chunkStatusHandler(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("./static/chunkstatus.html")
	if err != nil {
		log.Println("Error serving HTML file:", err)
		http.Error(w, "Internal Server Error, HTML File Not Found", http.StatusInternalServerError)
	} else {
		http.ServeFile(w, r, "./static/chunkstatus.html")
	}
}
