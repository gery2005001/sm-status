package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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

	//过滤
	smid := &pbv2alpha1.RewardStreamRequest_Smesher{
		Smesher: s,
	}

	// 获取reward
	reqReward := &pbv2alpha1.RewardStreamRequest{
		StartLayer: l,
		EndLayer:   l,
		FilterBy:   smid,
	}
	resReward, err := client.Stream(ctx, reqReward)
	if err != nil {
		log.Printf("get reward stream error: %s\n", err.Error())
		return 0, err
	}

	total := uint64(0)

	for {
		event, err := resReward.Recv()
		if err != nil {
			log.Println("Query reward end")
			break
		}

		log.Println(event.String())
		//total += event.GetV1().Total

		reward, err := extractAndConvertNumber(event.String())
		if err != nil {
			log.Println(err)
		} else {
			total += reward
		}

	}
	return total, nil
}

func extractAndConvertNumber(s string) (uint64, error) {
	// 分割字符串
	parts := strings.Split(s, " ")

	// 查找包含"3:"的部分
	var targetPart string
	found := false
	for _, part := range parts {
		if strings.HasPrefix(part, "3:") {
			targetPart = part
			found = true
			break
		}
	}

	// 如果没有找到"3:"，返回错误
	if !found {
		return 0, errors.New("未找到'3:'标记")
	}

	// 提取数字部分
	numberStr := strings.TrimPrefix(targetPart, "3:")

	// 检查提取的字符串是否为空
	if numberStr == "" {
		return 0, errors.New("'3:'后没有数字")
	}

	// 将字符串转换为uint64
	number, err := strconv.ParseUint(numberStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("转换为uint64失败: %v", err)
	}

	return number, nil
}
