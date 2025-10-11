package main

import (
	"MediaTools/internal/app"
	"MediaTools/internal/info"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	isDev       bool
	isServer    bool
	port        uint
	showVersion bool
)

func center(s string, width int, fill string) string {
	if len(s) >= width {
		return s
	}
	padding := width - len(s)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding
	return strings.Repeat(fill, leftPadding) + s + strings.Repeat(fill, rightPadding)
}

func init() {
	flag.BoolVar(&showVersion, "version", false, "显示版本信息\nShow version information")
	flag.BoolVar(&isDev, "dev", false, "是否启用开发者模式\nEnable developer mode")
	if info.Version.SupportDesktopMode { // 桌面模式下允许切换服务器模式
		flag.BoolVar(&isServer, "server", false, "是否启用 Web 服务器模式\nEnable web server mode")
	} else { // 服务器模式下强制启用服务器模式
		isServer = true
	}
	flag.UintVar(&port, "port", 5000, "Web 服务器端口（默认 5000）\nWeb server port (default 5000)")
	flag.Parse()

	fmt.Print("\033[2J") // 清屏
	fmt.Println(info.ProjectLogo)
	fmt.Println(center(
		fmt.Sprintf(" MediaWarp %s ", info.Version.AppVersion),
		81,
		"=",
	))
	gin.SetMode(gin.ReleaseMode)
}

// @title MediaTools API 文档
// @version 1.0
// @description 下一代媒体刮削&整理工具
// @Schemes HTTP
func main() {
	if showVersion {
		str, err := json.MarshalIndent(info.Version, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(str))
		return
	}

	if isServer {
		logrus.Infof("启动参数: 开发者模式=%v, 服务器模式=%v, 端口=%d", isDev, isServer, port)
	} else {
		port = 0 // 桌面模式下端口无效
		logrus.Infof("启动参数: 开发者模式=%v, 服务器模式=%v", isDev, isServer)
	}

	app.InitApp(isDev, isServer, port)
	defer logrus.Info("应用程序已退出")
	app.Run()
}
