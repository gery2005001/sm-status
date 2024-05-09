package utility

import (
	"fmt"
	"time"
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
