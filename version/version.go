package version

import "fmt"

var (
	// Version信息
	Version string
	// 编译日期
	BuildDate string
	// 编译时间
	BuildTime string
	// GO版本
	GO_Version string
)

func PrintCLIVersion() string {
	return fmt.Sprintf("Version: %s, build on %s %s, %s", Version, BuildDate, BuildTime, GO_Version)
}
