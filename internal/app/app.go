//go:build desktop
// +build desktop

package app

import (
	"context"
)

type App struct {
	ctx            context.Context
	systrayEndfunc func()
}

func NewApp() *App {
	return &App{}
}

const (
	ShowWindowsString = "显示窗口"
	ShowWindowsTip    = "显示应用窗口"
	HideWindowsString = "隐藏窗口"
	HideWindowsTip    = "隐藏应用窗口到系统托盘"
	QuitString        = "退出"
	QuitTip           = "退出应用"
)

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// switch version.Version.OS { // 设置托盘图标
	// case "darwin": // 支持 SVG 格式的系统
	// 	systray.SetIcon(web.GetLogoSVGData())
	// default:
	// 	systray.SetIcon(web.GetIconData())
	// }

	// mQuit := systray.AddMenuItem(QuitString, QuitTip)
	// mShowWindow := systray.AddMenuItem(ShowWindowsString, ShowWindowsTip)
	// mHideWindow := systray.AddMenuItem(HideWindowsString, HideWindowsTip)

	// go func() {
	// 	for {
	// 		select {
	// 		case <-mQuit.ClickedCh:
	// 			runtime.Quit(a.ctx)
	// 		case <-mShowWindow.ClickedCh:
	// 			runtime.Show(a.ctx)
	// 		case <-mHideWindow.ClickedCh:
	// 			runtime.Hide(a.ctx)
	// 		}
	// 	}
	// }()

	// startFunc, endFunc := systray.RunWithExternalLoop(nil, nil)
	// a.systrayEndfunc = endFunc
	// startFunc()
}

func (a *App) shutdown(ctx context.Context) {
	// if a.systrayEndfunc != nil {
	// 	a.systrayEndfunc()
	// }
}
