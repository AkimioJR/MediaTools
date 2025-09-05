package main

import (
	"MediaTools/internal/app"
	"MediaTools/internal/config"
	"embed"
	"flag"
	"fmt"
	"math/rand"
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
	isDev    bool
	isServer bool
	port     uint
)

func init() {
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
	for {
		port := rand.Intn(65535-20000) + 20000
		addr := fmt.Sprintf("localhost:%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port
		}
	}
}

// @title MediaTools API 文档
// @version 1.0
// @description 下一代媒体刮削&整理工具
// @Schemes HTTP
func main() {
	if !isServer {
		port = uint(findAvailablePort())
	}
	logrus.Infof("启动参数: 开发者模式=%v, 服务器模式=%v, 端口=%d", isDev, isServer, port)
	app.InitApp(isDev, isServer, port, &webDist)
	app.Run()
	logrus.Info("应用程序已退出")
}
