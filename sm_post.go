package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sm-status/utility"
	"sync"
	"time"
)

type Post struct {
	Enable          bool   `json:"enable"`
	Title           string `json:"title"`
	Capacity        string `json:"capacity"`
	OperatorAddress string `json:"operator-address"`
	Publish         SmPublish
	Eligs           []SmEligs
	SmesherId       []byte
	Status          string
	OaStatus        string
	NumUnits        uint32 `json:"numunits"`
}

type SmPublish struct {
	Time    time.Time
	Publish uint32
	Target  uint32
}

type SmEligs struct {
	Time  time.Time
	Epoch uint32
	Layer uint32
	Count uint32
	Total uint64
}

// Post相关函数

// 多线程获取Post Operator Status
func (x *Post) fetchPostOperator(wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()
	if !x.Enable {
		x.Status = ST_Disabled
		x.OaStatus = "Post disabled"
		//log.Println("post is disabled.")
		ch <- fmt.Sprintf("Post: %s OaStatus: %s", x.Title, x.OaStatus)
		return
	}
	if x.OperatorAddress == "" {
		x.Status = ST_Alone
		x.OaStatus = "No operator address"
		//log.Println("no set operator address.")
		ch <- fmt.Sprintf("Post: %s OaStatus: %s", x.Title, x.OaStatus)
		return
	}
	// 创建一个带有超时设置的 HTTP 客户端
	client := &http.Client{
		Timeout: GetTimeout() * time.Second, // 设置超时时间
	}

	// 创建一个上下文对象
	ctx, cancel := context.WithTimeout(context.Background(), GetTimeout()*time.Second)
	defer cancel() // 一定要确保取消上下文

	// log.Println("get status from: ", x.OperatorAddress)
	// 发送带有上下文的 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", x.OperatorAddress, nil)
	if err != nil {
		//log.Println("request operator failed:", err)
		x.Status = ST_Failed
		x.OaStatus = utility.GetHttpStatusCode(err)
		ch <- fmt.Sprintf("Post: %s OaStatus: %s", x.Title, x.OaStatus)
		return
	}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		//log.Println("request operator failed:", err)
		x.Status = ST_Failed
		x.OaStatus = utility.GetHttpStatusCode(err)
		ch <- fmt.Sprintf("Post: %s OaStatus: %s", x.Title, x.OaStatus)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//log.Println("operator read body failed:", err)
		x.Status = ST_Failed
		x.OaStatus = utility.GetHttpStatusCode(err)
		ch <- fmt.Sprintf("Post: %s OaStatus: %s", x.Title, x.OaStatus)
		return
	}

	x.Status = ST_Success
	x.OaStatus = string(body)
	ch <- fmt.Sprintf("Post: %s OaStatus: %s", x.Title, x.OaStatus)
	//log.Println("successfully get operator: ", x.OaStatus)
}

func (x *Post) GetStatusColorCSS() string {
	switch x.Status {
	case ST_Empty:
		return ST_Empty_CSS
	case ST_Alone:
		return ST_Alone_CSS
	case ST_Disabled:
		return ST_Disabled_CSS
	case ST_Running:
		return ST_Running_CSS
	case ST_Success:
		return ST_Success_CSS
	case ST_Failed:
		return ST_Failed_CSS
	default:
		return ST_Running_CSS
	}
}
