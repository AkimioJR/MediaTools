package main

import (
	"MediaTools/internal/app"
	"MediaTools/internal/version"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//go:embed web/dist
var webDist embed.FS

const LOGO = `
███╗   ███╗███████╗██████╗ ██╗ █████╗ ████████╗ ██████╗  ██████╗ ██╗     ███████╗
████╗ ████║██╔════╝██╔══██╗██║██╔══██╗╚══██╔══╝██╔═══██╗██╔═══██╗██║     ██╔════╝
██╔████╔██║█████╗  ██║  ██║██║███████║   ██║   ██║   ██║██║   ██║██║     ███████╗
██║╚██╔╝██║██╔══╝  ██║  ██║██║██╔══██║   ██║   ██║   ██║██║   ██║██║     ╚════██║
██║ ╚═╝ ██║███████╗██████╔╝██║██║  ██║   ██║   ╚██████╔╝╚██████╔╝███████╗███████║
╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝  ╚═════╝ ╚══════╝╚══════╝`

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
	if app.SupportDesktopMode { // 桌面模式下允许切换服务器模式
		flag.BoolVar(&isServer, "server", false, "是否启用 Web 服务器模式\nEnable web server mode")
	} else { // 服务器模式下强制启用服务器模式
		isServer = true
	}
	flag.UintVar(&port, "port", 5000, "Web 服务器端口（默认 5000）\nWeb server port (default 5000)")
	flag.Parse()

	fmt.Print("\033[2J") // 清屏
	fmt.Println(LOGO)
	fmt.Println(center(
		fmt.Sprintf(" MediaWarp %s ", version.Version.AppVersion),
		81,
		"=",
	))
	gin.SetMode(gin.ReleaseMode)
}

// findAvailablePort 查找一个可用的高位端口
func findAvailablePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// @title MediaTools API 文档
// @version 1.0
// @description 下一代媒体刮削&整理工具
// @Schemes HTTP
func main() {
	if showVersion {
		str, err := json.MarshalIndent(version.Version, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(str))
		return
	}
	if !isServer {
		port = uint(findAvailablePort())
	}
	logrus.Infof("启动参数: 开发者模式=%v, 服务器模式=%v, 端口=%d", isDev, isServer, port)
	app.InitApp(isDev, isServer, port, &webDist)
	app.Run()
	logrus.Info("应用程序已退出")
}
