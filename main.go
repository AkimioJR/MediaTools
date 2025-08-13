package main

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller"
	"MediaTools/internal/logging"
	"MediaTools/internal/router"
	"flag"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const LOGO = `
███╗   ███╗███████╗██████╗ ██╗ █████╗ ████████╗ ██████╗  ██████╗ ██╗     ███████╗
████╗ ████║██╔════╝██╔══██╗██║██╔══██╗╚══██╔══╝██╔═══██╗██╔═══██╗██║     ██╔════╝
██╔████╔██║█████╗  ██║  ██║██║███████║   ██║   ██║   ██║██║   ██║██║     ███████╗
██║╚██╔╝██║██╔══╝  ██║  ██║██║██╔══██║   ██║   ██║   ██║██║   ██║██║     ╚════██║
██║ ╚═╝ ██║███████╗██████╔╝██║██║  ██║   ██║   ╚██████╔╝╚██████╔╝███████╗███████║
╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝  ╚═════╝ ╚══════╝╚══════╝`

var (
	isDev bool
)

func init() {
	flag.BoolVar(&isDev, "dev", false, "是否启用开发者模式\nEnable developer mode")
	flag.Parse()

	fmt.Print("\033[2J") // 清屏
	fmt.Println(LOGO)
	fmt.Println(strings.Repeat("=", 31) + fmt.Sprintf(" MediaWarp %s ", config.Version.AppVersion) + strings.Repeat("=", 32))
	gin.SetMode(gin.ReleaseMode)
}

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

	logrus.Info("开始初始化全部工具链...")
	err = controller.InitAllControllers()
	if err != nil {
		panic(fmt.Sprintf("工具链初始化失败: %v", err))
	}
	logrus.Info("全部工具链初始化完成")

	ginR := router.InitRouter()
	err = ginR.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("启动服务器失败: %v", err))
	}
	logrus.Info("服务器启动成功，监听端口: 8080")
}
