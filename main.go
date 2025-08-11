package main

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller"
	"MediaTools/internal/logging"
	"MediaTools/internal/router"
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
╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝  ╚═════╝ ╚══════╝╚══════╝
                                                                                 `

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
	gin.SetMode(gin.ReleaseMode)
	fmt.Println(LOGO)
	fmt.Println(center(fmt.Sprintf(" MediaWarp %s ", config.Version.AppVersion), 81, "="))
}

func main() {
	logrus.Info("开始初始化配置...")
	err := config.Init()
	if err != nil {
		panic(err)
	}
	logrus.Info("配置初始化完成")

	logrus.Info("开始初始化日志...")
	logging.Init()
	logrus.Info("日志初始化完成")

	logrus.Info("开始初始化全部工具链...")
	err = controller.InitAllControllers()
	if err != nil {
		panic(err)
	}
	logrus.Info("全部工具链初始化完成")

	ginR := router.InitRouter()
	err = ginR.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("启动服务器失败: %v", err))
	}
	logrus.Info("服务器启动成功，监听端口: 8080")
}
