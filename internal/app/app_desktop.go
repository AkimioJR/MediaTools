//go:build !onlyServer

package app

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	webview "github.com/webview/webview_go"
)

const AppName = "MediaTools"

var SupportDesktopMode = true

func openWindows() {
	w := webview.New(false)
	defer w.Destroy()
	w.SetSize(800, 600, webview.HintNone)
	w.SetTitle(AppName)
	w.Navigate(fmt.Sprintf("http://localhost:%d", port))
	w.Run()
}

func runDesktop() {
	logrus.Infof("桌面模式启动中，端口: %d", port)
	go func() {
		runServer()
		os.Exit(0)
	}()
	openWindows()
}

func Run() {
	if isServer { // 启动服务器模式
		runServer()
	} else { // 启动桌面模式
		runDesktop()
	}
}
