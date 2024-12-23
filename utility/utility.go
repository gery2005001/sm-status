package utility

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 将传入的分钟数转换为格式为"2d 1h 0m"的字符串
func DurationToTimeFormat(d time.Duration) string {
	//log.Println("LeftTime:", d)
	tag := ""
	if d < 0 {
		tag = "-"
		d = -d
	}
	// 分别计算天、小时、分钟的值
	days := d / (24 * time.Hour)
	hours := (d % (24 * time.Hour)) / time.Hour
	minutes := (d % time.Hour) / time.Minute

	// 构造格式化后的字符串
	var result string

	result += fmt.Sprintf("%dd ", days)

	result += fmt.Sprintf("%dh ", hours)

	result += fmt.Sprintf("%dm", minutes)

	return tag + result
}

// 返回units的TB
func UnitsToTB(units uint32) string {
	if units <= 0 {
		return "0 TB"
	}
	result := float64(units * 64 * 1024)

	tb := result / (1024 * 1024)

	return fmt.Sprintf("%.2fTB", tb)
}

// 给一个时间戳添加或减去指定的秒
func TimeStampAddSecond(t uint32, s time.Duration) int64 {
	timeStamp := int64(t)
	//timeStamp转负成time.Time类型
	srcTime := time.Unix(timeStamp, 0)
	//增加指定的秒数
	newTime := srcTime.Add(s * time.Second)
	//获取新的时间戳
	newTimeStamp := newTime.Unix()

	return newTimeStamp
}

// 判断grpc服务错误码
func GetGRPCStatusCode(err error) string {
	if err != nil {
		// 获取错误状态码
		st := status.Code(err)

		// 判断不同类型的错误
		switch st {
		case codes.DeadlineExceeded:
			// 处理超时错误
			return "请求超时"
		case codes.Unavailable:
			// 处理服务不可用错误
			return "服务不可用"
		case codes.Internal:
			// 处理内部错误
			return "服务器内部错误"
		case codes.NotFound:
			// 处理未找到资源错误
			return "资源未找到"
		case codes.InvalidArgument:
			// 处理参数错误
			return "参数无效"
		default:
			// 处理其他类型错误
			return fmt.Sprintf("其他错误: %v", err)
		}
	}

	return "Successfully"
}

// 判断http服务错误码
func GetHttpStatusCode(err error) string {
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return "请求超时"
		case errors.Is(err, context.Canceled):
			return "请求被取消"
		case errors.Is(err, io.EOF):
			return "连接意外关闭"
		default:
			if _, ok := err.(net.Error); ok {
				return "网络连接超时"
			}

			var urlErr *url.Error
			if errors.As(err, &urlErr) {
				if urlErr.Timeout() {
					return "URL请求超时"
				}
				return fmt.Sprintf("URL错误: %v", urlErr)
			}

			return fmt.Sprintf("其他错误: %v\n", err)
		}
	}
	return "Successfully"
}
