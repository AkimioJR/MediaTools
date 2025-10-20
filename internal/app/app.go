//go:build desktop
// +build desktop

package app

import (
	"MediaTools/internal/info"
	"context"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func (app *App) newMenu() *menu.Menu {
	appMenu := menu.NewMenu()
	switch info.Version.OS {
	case "darwin":

		mainMenu := appMenu.AddSubmenu(info.ProjectName)
		mainMenu.AddText("关于", nil, func(_ *menu.CallbackData) {
			app.OpenAboutDialog()
		})

		mainMenu.AddSeparator()
		mainMenu.AddText("隐藏窗口", keys.CmdOrCtrl("h"), func(_ *menu.CallbackData) {
			app.Hide()
		})

		mainMenu.AddSeparator()
		mainMenu.AddText("退出", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
			app.Quit()
		})
	}

	return appMenu
}

func (app *App) Quit() {
	runtime.Quit(app.ctx)
}

func (app *App) Hide() {
	runtime.Hide(app.ctx)
}

func (app *App) Show() {
	runtime.Show(app.ctx)
}

const (
	closeBtnText = "关闭"
	jumpBtnText  = "跳转至项目主页"
)

func (app *App) OpenAboutDialog() {
	btnText, err := runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "关于 " + info.ProjectName,
		Message:       "一个用于媒体文件管理和处理的工具。\n\n" + info.Copyright + "\n\n" + info.Version.String(),
		Buttons:       []string{closeBtnText, jumpBtnText},
		CancelButton:  closeBtnText,
		DefaultButton: jumpBtnText,
	})
	if err != nil {
		return
	}
	if btnText == jumpBtnText {
		runtime.BrowserOpenURL(app.ctx, info.ProjectURL)
	}
}
