//go:build !onlyServer

package app

import (
	"MediaTools/internal/router"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	webview "github.com/webview/webview_go"
)

var SupportDesktopMode = true

func Run() {
	ginR := router.InitRouter(isDev, webDist)

	if isServer { // 启动服务器模式
		runServer()
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
