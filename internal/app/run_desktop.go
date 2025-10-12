//go:build desktop
// +build desktop

package app

import (
	"MediaTools/internal/info"
	"MediaTools/internal/router"
	"MediaTools/web"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func init() {
	info.Version.SupportDesktopMode = true
}

func newMenu(app *App) *menu.Menu {
	if info.Version.OS != "darwin" { // 非 Mac 系统不创建菜单
		return nil
	}

	appMenu := menu.NewMenu()

	mainMenu := appMenu.AddSubmenu("MediaTools")
	mainMenu.AddText("关于 "+info.ProjectName, nil, func(_ *menu.CallbackData) {
		runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
			Title:   "关于 " + info.ProjectName,
			Message: "一个用于媒体文件管理和处理的工具。\n\n" + info.Copyright + "\n\n" + info.Version.String(),
		})
	})
	mainMenu.AddSeparator()
	mainMenu.AddText("隐藏窗口", nil, func(_ *menu.CallbackData) {
		runtime.Hide(app.ctx)
	})

	mainMenu.AddSeparator()
	mainMenu.AddText("退出", nil, func(_ *menu.CallbackData) {
		runtime.Quit(app.ctx)
	})

	return appMenu
}

func runDesktop() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             info.ProjectName,
		Width:             1024,
		Height:            768,
		Menu:              newMenu(app),
		Frameless:         true, // 无边框窗口
		HideWindowOnClose: true, // 关闭窗口时隐藏应用
		OnStartup:         app.startup,
		OnShutdown:        app.shutdown,
		AssetServer: &assetserver.Options{
			Assets:  web.WebDist,
			Handler: router.InitRouter(info.RuntimeAppStatus.IsDev, nil),
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			BackdropType:         windows.Mica, // 使用Mica效果
			Theme:                windows.SystemDefault,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		panic("Error: " + err.Error())
	}
	logrus.Info("程序正常退出")
}

func Run() {
	if info.RuntimeAppStatus.DesktopMode { // 启动桌面模式
		runDesktop()
	} else { // 启动服务器模式
		runServer()
	}
}
