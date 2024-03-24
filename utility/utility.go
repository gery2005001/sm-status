package utility

import (
	"fmt"
	"time"
)

// 将传入的分钟数转换为格式为"2d 1h 0m"的字符串
func DurationToTimeFormat(d time.Duration) string {
	//log.Println("LeftTime:", d)
	if d < 0 {
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

	return result
}
