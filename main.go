package main

import (
	"MediaTools/internal/app"
	"MediaTools/internal/config"
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

func init() {
	flag.BoolVar(&showVersion, "version", false, "显示版本信息\nShow version information")
	flag.BoolVar(&isDev, "dev", false, "是否启用开发者模式\nEnable developer mode")
	flag.BoolVar(&isServer, "server", false, "是否启用 Web 服务器模式\nEnable web server mode")
	flag.UintVar(&port, "port", 8080, "Web 服务器端口（默认 8080）\nWeb server port")
	flag.Parse()

	fmt.Print("\033[2J") // 清屏
	fmt.Println(LOGO)
	fmt.Println(strings.Repeat("=", 31) + fmt.Sprintf(" MediaWarp %s ", config.Version.AppVersion) + strings.Repeat("=", 32))
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
		str, err := json.MarshalIndent(config.Version, "", "  ")
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
