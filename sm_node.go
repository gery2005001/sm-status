package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	pb "sm-status/spacemesh/v1"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ST_Empty            = ""
	ST_Running          = "R"
	ST_Failed           = "F"
	ST_Success          = "S"
	ST_Disabled         = "D"
	ST_Alone            = "A"
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
	Status              string
	PostInfo            []Post
	NodeVer             string
	Epoch               uint32
	SLayer              uint32
	TLayer              uint32
	VLayer              uint32
	IsSynced            bool
	Peers               uint64
	HasNewVer           bool
}

// Node相关函数

func (x *Node) GetStatusColorCSS() string {
	switch x.Status {
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

// 获取当前Epoch
func (x *Node) GetCurrentEpoch() {
	timeOut := GetTimeout()

	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPublicListener)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(timeOut))
	if err != nil {
		log.Println("Failed to connect: ", err)
		return
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
		log.Printf("Failed to call service: %v", err)
		return
	}

	x.Epoch = resEpoch.Epochnum.Number

}

// 从Node的GRPC服务中获取Node的状态
func (x *Node) GetNodeStatus() {
	timeout := GetTimeout()
	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPublicListener)

	log.Println("Starting get node status from ", grpcAddr)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(timeout))
	if err != nil {
		log.Println("Connect node failed: ", err)
		x.Status = ST_Failed
		return
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
		x.Status = ST_Failed
		return
	}
	x.NodeVer = resVer.VersionString.Value
	latestVer := GetLatestVer()
	if resVer.VersionString.Value != latestVer {
		x.HasNewVer = true
	}

	// 获取NodeStatus
	reqStatus := &pb.StatusRequest{}

	// 调用 gRPC 服务
	resStatus, err := client.Status(ctx, reqStatus)

	if err != nil {
		log.Printf("Failed to call service: %v", err)
		x.Status = ST_Failed
		return
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

	log.Println("Successfully get node status from ", grpcAddr)
}

// 获取Node所有Post的OperatorStatus
func (x *Node) GetPostOperatorStatus() {
	for i := 0; i < len(x.Post); i++ {
		x.Post[i].GetOperatorAddressStatus()
	}
}

// 从Node的GRPC服务中获取Events
func (x *Node) getEventsStreams() {
	timeout := GetTimeout()
	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPrivateListener)

	log.Println("Starting get Events Stream from ", grpcAddr)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(timeout))
	if err != nil {
		log.Println("Connect node failed: ", err)
		x.Status = ST_Failed
		return
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
		log.Printf("Failed to call service: %v", err)
		return
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
		newEvent := true
		//如果是 Eligibilities Event
		if _, ok := nd.(*pb.Event_Eligibilities); ok && newEvent {
			log.Printf("Found Eligs for %x \n", nEvent.GetEligibilities().Smesher)
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
			newEvent = false
		}
		//如果是 poetWaitProof Event
		if _, ok := nd.(*pb.Event_PoetWaitProof); ok && newEvent {
			log.Printf("Found poetWaitProof for %x", nEvent.GetPoetWaitProof().Smesher)
			for i, sm := range x.PostInfo {
				if bytes.Equal(sm.SmesherId, nEvent.GetPoetWaitProof().Smesher) {
					x.PostInfo[i].Publish.Time = nEvent.Timestamp.AsTime()
					x.PostInfo[i].Publish.Publish = nEvent.GetPoetWaitProof().Publish
					x.PostInfo[i].Publish.Target = nEvent.GetPoetWaitProof().Target
				}
			}
			newEvent = false
		}
	}

	log.Println("Successfully get Events Stream from ", grpcAddr)
}

// // 从Node的PostInfoService中获取PostInfo
func (x *Node) getPostInfoFromGRPC() {
	if x.NodeType != "multi" {
		x.PostInfo = nil
		alonePost := Post{
			Title:  x.Name,
			Status: x.NodeType,
		}
		x.PostInfo = append(x.PostInfo, alonePost)
		return
	}
	timeout := GetTimeout()
	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPostListener)

	log.Println("Starting get Post Info from ", grpcAddr)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(timeout))
	if err != nil {
		log.Println("Connect node failed: ", err)
		x.Status = ST_Failed
		return
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
		log.Printf("Failed to call service: %v", err)
		return
	}
	if len(response.States) > 0 {
		x.PostInfo = []Post{}
		for i := 0; i < len(response.States); i++ {
			x.PostInfo = append(x.PostInfo, Post{
				SmesherId: response.States[i].Id,
				Title:     response.States[i].Name,
				Status:    response.States[i].State.String(),
			})

		}
	}
	log.Println("Successfully get Post Info from ", grpcAddr)
}
