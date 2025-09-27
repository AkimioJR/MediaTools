//go:build !onlyServer

package app

import (
	"fmt"

	"fyne.io/systray"
	"github.com/sirupsen/logrus"
	webview "github.com/webview/webview_go"
)

const AppName = "MediaTools"

var (
	SupportDesktopMode = true
	globalView         = webview.New(false)
)

func showWindow() {
	logrus.Debug("显示窗口")
	globalView.Dispatch(func() {
		globalView.Navigate(fmt.Sprintf("http://localhost:%d", port))
		globalView.SetSize(800, 600, webview.HintNone)
	})
}

func hideWindow() {
	logrus.Debug("隐藏窗口")
	globalView.Dispatch(func() {
		globalView.SetHtml("")
		globalView.SetSize(0, 0, webview.HintNone)
	})
}

func onReady() {
	systray.SetTitle(AppName)
	systray.SetTooltip(AppName + " - 工具栏")
	quitItem := systray.AddMenuItem("退出", "退出应用程序")
	showWindowsItem := systray.AddMenuItem("显示窗口", "显示应用程序窗口")
	hideWindowsItem := systray.AddMenuItem("隐藏窗口", "隐藏应用程序窗口")
	go func() {
		for {
			select {
			case <-showWindowsItem.ClickedCh:
				showWindow()
			case <-hideWindowsItem.ClickedCh:
				hideWindow()
			case <-quitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()
	logrus.Debug("系统托盘启动成功")
}

func onExit(quitFlag *bool) func() {
	return func() {
		logrus.Info("退出应用程序中...")
		globalView.Terminate()
		*quitFlag = true
	}
}

func runDesktop() {
	logrus.Infof("桌面模式启动中，端口: %d", port)
	go func() {
		runServer()
		systray.Quit()
	}()

	quitFlag := false
	startFn, endFn := systray.RunWithExternalLoop(onReady, onExit(&quitFlag))
	startFn()
	defer endFn()

	globalView.SetTitle(AppName)
	defer globalView.Destroy()
	for !quitFlag {
		showWindow()
		logrus.Debug("webview 运行中...")
		globalView.Run()
		logrus.Debug("webview 运行结束")
		systray.Quit()
	}

}

func Run() {
	if isServer { // 启动服务器模式
		runServer()
	} else { // 启动桌面模式
		runDesktop()
	}
}
