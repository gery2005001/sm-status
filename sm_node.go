package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	pb "sm-status/spacemesh/v1"
	"sm-status/utility"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ST_Empty            = ""  //未初始化
	ST_Running          = "R" //运行中未同步
	ST_Failed           = "F" //通讯失败，未开机
	ST_Success          = "S" //开机并同步
	ST_Disabled         = "D" //配置未开启
	ST_Alone            = "A" //单机节点
	ST_Empty_CSS        = "st-running"
	ST_Failed_CSS       = "st-failed"
	ST_Success_CSS      = "st-success"
	ST_Disabled_CSS     = "st-disabled"
	ST_Alone_CSS        = "st-alone"
	ST_Running_CSS      = "st-running"
	SM_LayerDuration    = 300
	SM_GetNewVerAddress = "https://api.github.com/repos/spacemeshos/go-spacemesh/releases/latest"
)

type Node struct {
	Name                string `json:"name"`
	IP                  string `json:"ip"`
	GrpcPublicListener  int    `json:"grpc-public-listener"`
	GrpcPrivateListener int    `json:"grpc-private-listener"`
	GrpcPostListener    int    `json:"grpc-post-listener"`
	GrpcJsonListener    int    `json:"grpc-json-listener"`
	Enable              bool   `json:"enable"`
	NodeType            string `json:"node-type"`
	Post                []Post `json:"post"`
	PostInfo            []Post
	NodeVer             string
	Epoch               uint32
	SLayer              uint32
	TLayer              uint32
	VLayer              uint32
	IsSynced            bool
	Peers               uint64
	HasNewVer           bool
	Status              string
}

// Node相关函数
func (x *Node) GetStatusColorCSS() string {
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
	default:
		return ST_Failed_CSS
	}
}

// 获取Node所有信息
func (x *Node) GetAllInformation() {
	//从node获取当前Epoch
	if err := x.getCurrentEpoch(); err != nil {
		x.setNodeToFailedStatus()
		log.Println(err)
		return
	}

	//从node获取version和status
	if err := x.getNodeVerAndStatus(); err != nil {
		x.setNodeToFailedStatus()
		log.Println(err)
		return
	}

	//从node的PostService中获取
	if err := x.getPostInfoState(); err != nil {
		log.Println(err)
		return
	}
}

// 清除node相关状态信息
func (x *Node) setNodeToFailedStatus() {
	x.Status = ST_Failed
	x.Epoch = 0
	x.SLayer = 0
	x.TLayer = 0
	x.VLayer = 0
	x.Peers = 0
	x.IsSynced = false
	x.NodeVer = ""
	x.PostInfo = []Post{}
}

// 从node获取当前Epoch，用以判断node是否开启
func (x *Node) getCurrentEpoch() error {
	if !x.Enable {
		return fmt.Errorf("node %s ip %s is disabled", x.Name, x.IP)
	}
	timeOut := GetTimeout()

	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPublicListener)

	log.Println("Starting get node current epoch from ", grpcAddr)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(timeOut))
	if err != nil {
		log.Println("get current epoch failed.")
		return err
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pb.NewMeshServiceClient(conn)

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), timeOut*time.Second)
	defer cancel()

	// 查询当前Epoch
	reqEpoch := pb.CurrentEpochRequest{}

	resEpoch, err := client.CurrentEpoch(ctx, &reqEpoch)
	if err != nil {
		log.Println("get current eopch failed.")
		return err
	}
	x.Epoch = resEpoch.Epochnum.Number

	x.Status = ST_Running
	log.Println("successfully get node current epoch from ", grpcAddr)
	return nil
}

// 从Node的GRPC服务中获取Node的version和status
func (x *Node) getNodeVerAndStatus() error {
	timeout := GetTimeout()
	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPublicListener)

	log.Println("starting get node version and status from ", grpcAddr)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(timeout))
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pb.NewNodeServiceClient(conn)

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// 获取node的版本号
	reqVer := &emptypb.Empty{}
	resVer, err := client.Version(ctx, reqVer)
	if err != nil {
		return err
	}
	x.NodeVer = resVer.VersionString.Value

	//获取客户端最新版本
	latestVer := GetLatestVer()
	if resVer.VersionString.Value != latestVer {
		x.HasNewVer = true
	} else {
		x.HasNewVer = false
	}

	// 获取NodeStatus
	reqStatus := &pb.StatusRequest{}

	// 调用 gRPC 服务
	resStatus, err := client.Status(ctx, reqStatus)

	if err != nil {
		return err
	}

	x.IsSynced = resStatus.Status.IsSynced
	if x.IsSynced {
		x.Status = ST_Success
	} else {
		x.Status = ST_Running
	}
	x.Peers = resStatus.Status.ConnectedPeers
	x.SLayer = resStatus.Status.SyncedLayer.Number
	x.TLayer = resStatus.Status.TopLayer.Number
	x.VLayer = resStatus.Status.VerifiedLayer.Number

	log.Println("successfully get node version and status from ", grpcAddr)
	return nil
}

// // 从Node的PostInfoService中获取PostInfo
func (x *Node) getPostInfoState() error {
	x.PostInfo = []Post{}
	if x.NodeType != "multi" {
		alonePost := Post{
			Title:  x.Name,
			Status: x.NodeType,
		}
		x.PostInfo = append(x.PostInfo, alonePost)
		return nil
	}
	timeout := GetTimeout()
	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPostListener)

	log.Println("starting get Post Info from ", grpcAddr)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(timeout))
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pb.NewPostInfoServiceClient(conn)

	// 设置超时时间为 3 秒
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// 构造请求
	request := &pb.PostStatesRequest{}

	// 调用 gRPC 服务
	response, err := client.PostStates(ctx, request)
	if err != nil {
		return err
	}
	if len(response.States) > 0 {
		for i := 0; i < len(response.States); i++ {
			x.PostInfo = append(x.PostInfo, Post{
				SmesherId: response.States[i].Id,
				Title:     response.States[i].Name,
				Status:    response.States[i].State.String(),
			})

		}
	}
	log.Println("successfully get Post Info from ", grpcAddr)

	if err := x.getEventsStreams(); err != nil {
		return err
	}

	return nil
}

// 从Node的GRPC服务中获取Events
func (x *Node) getEventsStreams() error {
	if x.NodeType == "smapp" {
		log.Println("node is smapp skip get events stream")
		return nil
	}
	timeout := GetTimeout()
	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPrivateListener)

	log.Println("starting get Events Stream from ", grpcAddr)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(timeout))
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pb.NewAdminServiceClient(conn)

	// 设置超时时间为 3 秒
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	// 构造请求
	request := &pb.EventStreamRequest{}

	// 调用 gRPC 服务
	esClient, err := client.EventsStream(ctx, request)
	if err != nil {
		return err
	}

	nEvent := &pb.Event{}
	for {
		nEvent, err = esClient.Recv()
		if err != nil {
			log.Println("Query ends")
			break
		}
		//获取事件类型
		nd := nEvent.Details
		//是否已处理
		isNew := true
		//如果是 Eligibilities Event
		if _, ok := nd.(*pb.Event_Eligibilities); ok && isNew {
			log.Printf("found eligs for %x \n", nEvent.GetEligibilities().Smesher)
			for i, sm := range x.PostInfo {
				if bytes.Equal(sm.SmesherId, nEvent.GetEligibilities().Smesher) {
					tmElgs := nEvent.GetEligibilities()
					for _, elg := range tmElgs.Eligibilities {
						x.PostInfo[i].Eligs = append(x.PostInfo[i].Eligs, SmEligs{
							Time:  nEvent.Timestamp.AsTime(),
							Epoch: nEvent.GetEligibilities().GetEpoch(),
							Layer: elg.Layer,
							Count: elg.Count,
						})
					}
				}
			}
			isNew = false
		}
		//如果是 poetWaitProof Event
		if _, ok := nd.(*pb.Event_PoetWaitProof); ok && isNew {
			log.Printf("found poetWaitProof for %x", nEvent.GetPoetWaitProof().Smesher)
			for i, sm := range x.PostInfo {
				if bytes.Equal(sm.SmesherId, nEvent.GetPoetWaitProof().Smesher) {
					x.PostInfo[i].Publish.Time = nEvent.Timestamp.AsTime()
					x.PostInfo[i].Publish.Publish = nEvent.GetPoetWaitProof().Publish
					x.PostInfo[i].Publish.Target = nEvent.GetPoetWaitProof().Target
				}
			}
			isNew = false
		}
	}
	log.Println("successfully get Events Stream from ", grpcAddr)
	return nil
}

// 获取node的状态html
func (x *Node) GetNodeStatusTableHTMLString() string {
	if !x.Enable {
		return ""
	}
	//获取node状态
	nodeSyncedText := ""
	nodeSTColor := x.GetStatusColorCSS()
	if x.Status == ST_Success && x.IsSynced {
		nodeSyncedText = "【已同步】"
	} else {
		if x.Status == ST_Empty {
			nodeSyncedText = "【获取中】"
		} else {
			nodeSyncedText = "【未同步】"
		}
		nodeSTColor = x.GetStatusColorCSS()
	}
	verSTColor := ""
	if x.HasNewVer {
		verSTColor = ST_Failed_CSS
	}
	//生成页面
	htmlData := "<table>"
	htmlData += "<colgroup><col class=\"col-per-15\"><col class=\"col-per-20\"><col class=\"col-per-15\"><col classe=\"col-per-10\"><col classe=\"auto-column\"></colgroup>"
	htmlData += "<thead>"
	htmlData += "<tr class=\"node-info\"><td class=\"td-left\" colspan=\"5\">"
	htmlData += fmt.Sprintf("<span>状态：<b>"+"<span class=\"%s\">%s</span></b></span>", nodeSTColor, nodeSyncedText)
	htmlData += "<span>　Node名称：<b>" + x.Name + "</b></span>　<span>IP：<b>" + x.IP + "</b></span>"
	htmlData += fmt.Sprintf("<span>　版本：<span class=\"%s\"><b>%s</b></span></span>", verSTColor, x.NodeVer)
	htmlData += fmt.Sprintf("　<span><span>Peers：<b>%d</b></span>", x.Peers)
	htmlData += fmt.Sprintf("　<span>Synced Layer：<b>%d</b></span>", x.SLayer)
	htmlData += fmt.Sprintf("　<span>Top Layer：<b>%d</b></span>", x.TLayer)
	htmlData += fmt.Sprintf("　<span>Verified Layer：<b>%d</b></span>", x.VLayer)
	htmlData += fmt.Sprintf("　<span>Epoch：<b>%d</b></span>", x.Epoch)
	htmlData += "</td></tr>"
	if x.PostInfo != nil {
		htmlData += "<thead><tr><th>KEY</th><th>State</th><th>Eligibilities</th><th>Publish</th><th>ID</th></tr></thead>"
		htmlData += "<tbody>"
		for i := 0; i < len(x.PostInfo); i++ {
			elgMsg := ""
			elgBn := ""
			elgEnd := "✓"
			leftTime := ""
			elgBtnStyle := "btn-running"
			for _, elg := range x.PostInfo[i].Eligs {
				if elg.Epoch >= x.Epoch {
					if elg.Layer == x.TLayer {
						elgBtnStyle = "btn-running"
						elgEnd = "【now】"
					} else if elg.Layer < x.TLayer {
						lt := (x.TLayer - elg.Layer) * SM_LayerDuration
						elgBtnStyle = "btn-success"
						leftTime = "-" + utility.DurationToTimeFormat(time.Duration(lt)*time.Second)
						elgEnd = "【✓】"
					} else {
						lt := (elg.Layer - x.TLayer) * SM_LayerDuration
						leftTime = utility.DurationToTimeFormat(time.Duration(lt) * time.Second)
						elgEnd = fmt.Sprintf("%s【%d】", leftTime, elg.Count)
					}
					//elgMsg = fmt.Sprintf("<span class=\"%s\">【%s】</span>Layer:<b>%d</b>,Count:%d", bkColor, leftTime, elg.Layer, elg.Count)
					elgMsg = fmt.Sprintf("【%s】Epoch:【%d】,Layer:【%d】,Count:【%d】", leftTime, elg.Epoch, elg.Layer, elg.Count)
					elgBn = fmt.Sprintf("<button class=\"%s\" onclick=\"alert('%s')\">%s</button>", elgBtnStyle, elgMsg, elgEnd)
				} else {
					elgBn = ""
				}
			}
			pwpMsg := ""
			pwpBn := ""
			if x.PostInfo[i].Publish.Publish >= x.Epoch {
				pwpMsg = fmt.Sprintf("Publish:【%d】,Target:【%d】", x.PostInfo[i].Publish.Publish, x.PostInfo[i].Publish.Target)
				pwpBn = fmt.Sprintf("<button class=\"btn-success\" onclick=\"alert('%s')\">【%d】</button>", pwpMsg, x.PostInfo[i].Publish.Target)
			} else {
				pwpBn = ""
			}
			htmlData += fmt.Sprintf("<tr><td>%s</td><td  class=\"td-left\">%s</td><td>%s</td><td>%s</td><td class=\"td-rtl\">%x</td><tr>", x.PostInfo[i].Title, x.PostInfo[i].Status, elgBn, pwpBn, x.PostInfo[i].SmesherId)
		}
		htmlData += "</tbody>"
	}

	htmlData += "</table>"
	return htmlData
}

// 获取Node所有Post的OperatorStatus
func (x *Node) getNodePostOperatorStatus() {
	if !x.Enable || x.Status == ST_Failed {
		log.Println("node is disabled or failed skip get Operator status")
		for i := 0; i < len(x.Post); i++ {
			x.Post[i].Status = ST_Failed
		}
		return
	}
	for i := 0; i < len(x.Post); i++ {
		x.Post[i].getPostOperator()
	}
}
