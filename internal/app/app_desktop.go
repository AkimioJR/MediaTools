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
	globalView         webview.WebView
	isWindowVisible    = false
	showWindowChan     = make(chan struct{}, 1)
	updateTrayMenuChan = make(chan struct{}, 5) // 新增：用于更新托盘菜单状态
	quitFlag           = false
)

func createWebView() {
	if globalView != nil {
		globalView.Destroy()
	}
	globalView = webview.New(false)
	globalView.SetTitle(AppName)
}

func showWindow() {
	if !isWindowVisible {
		logrus.Debug("显示窗口")
		if globalView == nil {
			createWebView()
		}

		// 使用 Dispatch 确保在主线程中执行 UI 操作
		globalView.Dispatch(func() {
			globalView.SetSize(800, 600, webview.HintNone)
			globalView.Navigate(fmt.Sprintf("http://localhost:%d", port))
		})

		isWindowVisible = true

		// 通知等待线程窗口已显示
		select {
		case showWindowChan <- struct{}{}:
		default:
		}
	}
}

func onReady() {
	const (
		showTitle = "显示窗口"
		showTip   = "显示应用程序窗口"
		hideTitle = "隐藏窗口"
		hideTip   = "隐藏应用程序窗口"
	)
	systray.SetTitle(AppName)
	systray.SetTooltip(AppName + " - 工具栏")
	switchWindowStatusItem := systray.AddMenuItem(hideTitle, hideTip)
	quitItem := systray.AddMenuItem("退出", "退出应用程序")

	go func() {
		for {
			select {
			case <-switchWindowStatusItem.ClickedCh:
				if isWindowVisible {
					logrus.Debug("用户从托盘隐藏窗口")
					isWindowVisible = false
					// 安全地终止当前 webview
					if globalView != nil {
						globalView.Dispatch(func() {
							globalView.Terminate()
						})
					}
				} else {
					logrus.Debug("用户从托盘显示窗口")
					showWindow()
				}
				updateTrayMenuChan <- struct{}{}

			case <-updateTrayMenuChan:
				// 根据当前窗口状态更新托盘菜单
				if isWindowVisible {
					switchWindowStatusItem.SetTitle(hideTitle)
					switchWindowStatusItem.SetTooltip(hideTip)
				} else {
					switchWindowStatusItem.SetTitle(showTitle)
					switchWindowStatusItem.SetTooltip(showTip)
				}
				logrus.Debug("已更新系统托盘菜单状态")

			case <-quitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()
	logrus.Debug("系统托盘启动成功")
}

func onExit() func() {
	return func() {
		logrus.Info("退出应用程序中...")
		quitFlag = true
		if globalView != nil {
			// 确保在主线程中终止 webview
			globalView.Dispatch(func() {
				globalView.Terminate()
			})
		}
	}
}

func runDesktop() {
	logrus.Infof("桌面模式启动中，端口: %d", port)
	go func() {
		runServer()
		systray.Quit()
	}()

	startFn, endFn := systray.RunWithExternalLoop(onReady, onExit())
	startFn()
	defer endFn()

	// 标记是否首次启动
	firstLaunch := true

	// 主循环：持续运行直到用户明确选择退出
	for !quitFlag {
		// 只在首次启动或用户主动请求时显示窗口
		if firstLaunch {
			// 首次启动，创建并显示窗口
			createWebView()
			defer globalView.Destroy()
			showWindow()
			firstLaunch = false

			// 确保托盘菜单状态正确
			select {
			case updateTrayMenuChan <- struct{}{}:
			default:
			}
		} else {
			logrus.Debug("等待用户通过系统托盘重新打开窗口")
			<-showWindowChan // 等待用户通过系统托盘请求显示窗口

			if !quitFlag {
				createWebView() // 创建新的 webview 实例
				defer globalView.Destroy()
				// showWindow() 已经在托盘处理中被调用
			}
		}

		// 如果没有退出，运行 webview
		if !quitFlag {
			logrus.Debug("webview 运行中...")
			globalView.Run()
			logrus.Debug("webview 运行结束")

			// 当 webview.Run() 退出时（用户关闭窗口），更新状态
			isWindowVisible = false

			// 通知系统托盘更新菜单状态
			select {
			case updateTrayMenuChan <- struct{}{}:
			default:
			}
		}
	}
	logrus.Info("程序正常退出")
}

func Run() {
	if isServer { // 启动服务器模式
		runServer()
	} else { // 启动桌面模式
		runDesktop()
	}
}
