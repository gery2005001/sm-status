package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	pbv2alpha1 "github.com/spacemeshos/api/release/go/spacemesh/v2alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 从Node的GRPC服务中获取指定Layer中Smesher的Reward
func (x *Node) GetLayerRewardWithSmesher(l uint32, s []byte) (uint64, error) {
	grpcAddr := fmt.Sprintf("%s:%d", x.IP, x.GrpcPrivateListener)

	log.Printf("starting get smesher %x reward in layer %d \n", s, l)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Printf("get node %s version error: %s\n", x.Name, err.Error())
		return 0, err
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pbv2alpha1.NewRewardStreamServiceClient(conn)

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), GetTimeout()*time.Second)
	defer cancel()

	// 获取node的版本号
	reqReward := &pbv2alpha1.RewardStreamRequest{
		StartLayer: l,
		EndLayer:   l,
	}
	resReward, err := client.Stream(ctx, reqReward)
	if err != nil {
		log.Printf("get reward stream error: %s\n", err.Error())
		return 0, err
	}

	for {
		event, err := resReward.Recv()
		if err != nil {
			log.Println("Query reward end")
			break
		}

		if bytes.Equal(event.GetV1().Smesher, s) {
			log.Printf("layer: %d,reward: %d,smesher: %x \n", event.GetV1().Layer, event.GetV1().Total, event.GetV1().Smesher)
			return event.GetV1().Total, nil
		}
		//log.Println(event.GetV1().GetSmesher())

	}
	return 0, fmt.Errorf("reward not found")
}
