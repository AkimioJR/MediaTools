//go:build desktop
// +build desktop

package app

import (
	"MediaTools/web"
	_ "embed"
	"runtime"

	"fyne.io/systray"
	"github.com/sirupsen/logrus"
	webview "github.com/webview/webview_go"
)

const AppName = "MediaTools"

type UIManager struct {
	url                string
	view               webview.WebView
	isVisible          bool
	showVisibleChan    chan struct{}
	updateTrayMenuChan chan struct{}
	quitFlag           bool
}

func NewUIManager(url string) *UIManager {
	return &UIManager{
		url:                url,
		view:               webview.New(false),
		isVisible:          false,
		showVisibleChan:    make(chan struct{}, 1),
		updateTrayMenuChan: make(chan struct{}, 1),
		quitFlag:           false,
	}
}

func (ui *UIManager) createWebView() {
	if ui.view != nil {
		ui.view.Destroy()
		ui.view = nil
	}
	ui.view = webview.New(false)
	ui.view.SetTitle(AppName)
}

func (ui *UIManager) showWindow() {
	if !ui.isVisible {
		logrus.Debug("显示窗口")
		if ui.view == nil {
			ui.createWebView()
		}

		// 使用 Dispatch 确保在主线程中执行 UI 操作
		ui.view.Dispatch(func() {
			ui.view.SetSize(800, 600, webview.HintNone)
			ui.view.Navigate(ui.url)
		})

		ui.isVisible = true

		ui.showVisibleChan <- struct{}{} // 通知等待线程窗口已显示
	}
}

func (ui *UIManager) onReady() {
	const (
		showTitle = "显示窗口"
		showTip   = "显示应用程序窗口"
		hideTitle = "隐藏窗口"
		hideTip   = "隐藏应用程序窗口"
	)
	switch runtime.GOOS {
	case "darwin": // 支持 SVG 图标系统
		systray.SetIcon(web.GetLogoSVGData())
	default:
		systray.SetIcon(web.GetIconData())
	}
	systray.SetTooltip(AppName + " - 工具栏")
	switchWindowStatusItem := systray.AddMenuItem(hideTitle, hideTip)
	setTitle := func(isShown bool) {
		if isShown {
			switchWindowStatusItem.SetTitle(hideTitle)
			switchWindowStatusItem.SetTooltip(hideTip)
		} else {
			switchWindowStatusItem.SetTitle(showTitle)
			switchWindowStatusItem.SetTooltip(showTip)
		}
	}
	quitItem := systray.AddMenuItem("退出", "退出应用程序")

	go func() {
		for {
			select {
			case <-switchWindowStatusItem.ClickedCh:
				if ui.isVisible {
					logrus.Debug("用户从托盘隐藏窗口")
					ui.isVisible = false
					if ui.view != nil { // 安全地终止当前 webview
						ui.view.Terminate()
					}
					setTitle(false)
				} else {
					logrus.Debug("用户从托盘显示窗口")
					ui.showWindow()
					setTitle(true)
				}

			case <-ui.updateTrayMenuChan:
				setTitle(ui.isVisible)
				logrus.Debug("已更新系统托盘菜单状态")

			case <-quitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()
	logrus.Debug("系统托盘启动成功")
}

func (ui *UIManager) onExit() {
	logrus.Info("退出应用程序中...")
	ui.quitFlag = true
	if ui.view != nil {
		// 确保在主线程中终止 webview
		ui.view.Dispatch(func() {
			ui.view.Terminate()
		})
	}
}

func (ui *UIManager) Run() {
	startFn, endFn := systray.RunWithExternalLoop(ui.onReady, ui.onExit)
	startFn()
	defer endFn()

	firstLaunch := true
	for !ui.quitFlag { // 只在首次启动或用户主动请求时显示窗口
		if firstLaunch { // 首次启动，创建并显示窗口
			ui.createWebView()
			defer ui.view.Destroy()
			ui.showWindow()
			firstLaunch = false
		} else {
			logrus.Debug("等待用户通过系统托盘重新打开窗口")
			<-ui.showVisibleChan // 等待用户通过系统托盘请求显示窗口

			if !ui.quitFlag {
				ui.createWebView() // 创建新的 webview 实例
				defer ui.view.Destroy()
			}
		}

		if !ui.quitFlag { // 如果没有退出，运行 webview
			logrus.Debug("webview 运行中...")
			ui.view.Run()
			logrus.Debug("webview 运行结束")
			ui.isVisible = false                // 当 webview.Run() 退出时（用户关闭窗口），更新状态
			ui.updateTrayMenuChan <- struct{}{} // 通知系统托盘更新菜单状态
		}
	}
}
