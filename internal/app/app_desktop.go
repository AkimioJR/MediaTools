//go:build !onlyServer

package app

import (
	"fmt"

	"fyne.io/systray"
	"github.com/sirupsen/logrus"
)

var SupportDesktopMode = true

func runDesktop() {
	logrus.Infof("桌面模式启动中，端口: %d", port)
	go func() {
		runServer()
		systray.Quit()
	}()
	ui := NewUIManager(fmt.Sprintf("http://localhost:%d", port))
	ui.Run()
	logrus.Info("程序正常退出")
}

func Run() {
	if isServer { // 启动服务器模式
		runServer()
	} else { // 启动桌面模式
		runDesktop()
	}
}
