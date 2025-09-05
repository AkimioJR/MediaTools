package main

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller"
	"MediaTools/internal/database"
	"MediaTools/internal/logging"
	"MediaTools/internal/router"
	"embed"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/polevpn/webview"
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
	logrus.Info("开始初始化配置...")
	err := config.Init()
	if err != nil {
		panic(fmt.Sprintf("配置初始化失败: %v", err))
	}
	logrus.Info("配置初始化完成")

	logrus.Info("开始初始化日志...")
	err = logging.Init()
	if err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
	logrus.Info("日志初始化完成")
	err = database.Init()
	if err != nil {
		panic(fmt.Sprintf("数据库初始化失败: %v", err))
	}

	logrus.Info("开始初始化全部工具链...")
	err = controller.InitAllControllers()
	if err != nil {
		panic(fmt.Sprintf("工具链初始化失败: %v", err))
	}
	logrus.Info("全部工具链初始化完成")

	ginR := router.InitRouter(isDev, &webDist)

	if isServer { // 启动服务器模式
		err = ginR.Run(":" + strconv.Itoa(int(port)))
		if err != nil {
			panic(fmt.Sprintf("启动服务器失败: %v", err))
		}
		logrus.Infof("服务器启动成功，监听端口: %d", port)
	} else { // 启动桌面模式
		// 查找可用端口
		port := findAvailablePort()
		logrus.Infof("桌面模式启动中，端口: %d", port)

		// 创建 HTTP 服务器
		srv := &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", port),
			Handler: ginR,
		}

		// 在后台启动服务器
		go func() {
			logrus.Infof("服务器启动在端口: %d", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Errorf("服务器启动失败: %v", err)
			}
		}()

		w := webview.New(800, 600, false, true)
		defer w.Destroy()
		w.SetSize(800, 600, webview.HintNone)
		w.SetTitle("MediaTools")
		w.Navigate(fmt.Sprintf("http://localhost:%d", port))
		w.Run()

		logrus.Info("应用程序已退出")
	}

}
