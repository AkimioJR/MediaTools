package app

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller"
	"MediaTools/internal/database"
	"MediaTools/internal/logging"
	"MediaTools/internal/router"
	"embed"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	webview "github.com/webview/webview_go"
)

var (
	isDev    bool // 是否启用开发者模式
	isServer bool // 是否启用 Web 服务器模式
	port     uint // Web 服务器端口
	webDist  *embed.FS
)

func InitApp(d bool, s bool, p uint, w *embed.FS) {
	// 初始化全局变量
	isDev = d
	isServer = s
	port = p
	webDist = w

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

}

func Run() {
	ginR := router.InitRouter(isDev, webDist)

	if isServer { // 启动服务器模式
		err := ginR.Run(":" + strconv.Itoa(int(port)))
		if err != nil {
			panic(fmt.Sprintf("启动服务器失败: %v", err))
		}
		logrus.Infof("服务器启动成功，监听端口: %d", port)
	} else { // 启动桌面模式
		logrus.Infof("桌面模式启动中，端口: %d", port)

		// 在后台启动服务器
		go func() {
			err := ginR.Run("localhost:" + strconv.Itoa(int(port)))
			if err != nil {
				panic(fmt.Sprintf("启动服务器失败: %v", err))
			}
			logrus.Infof("服务器启动成功，监听端口: %d", port)
		}()

		w := webview.New(false)
		defer w.Destroy()
		w.SetSize(800, 600, webview.HintNone)
		w.SetTitle("MediaTools")
		w.Navigate(fmt.Sprintf("http://localhost:%d", port))
		w.Run()
	}
}
