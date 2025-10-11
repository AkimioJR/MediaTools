//go:build desktop
// +build desktop

package app

import (
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

const (
	Copyright = "Copyright © 2025 AKimioJR(akimio.jr@gmail.com)"
)

var SupportDesktopMode = true

func newMenu(app *App) *menu.Menu {
	appMenu := menu.NewMenu()
	// if version.Version.OS == "" {
	// 	appMenu.Append(menu.AppMenu())
	// } else {
	mainMenu := appMenu.AddSubmenu("MediaTools")
	mainMenu.AddText("关于 "+ProjectName, nil, func(_ *menu.CallbackData) {
		runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
			Title:   "关于 " + ProjectName,
			Message: "MediaTools 是一个用于媒体文件管理和处理的工具。\n" + Copyright,
		})
	})
	mainMenu.AddSeparator()
	mainMenu.AddText("隐藏窗口", nil, func(_ *menu.CallbackData) {
		runtime.Hide(app.ctx)
	})
	// mainMenu.AddText("显示窗口", nil, func(_ *menu.CallbackData) {
	// 	runtime.Show(app.ctx)
	// })
	mainMenu.AddSeparator()
	mainMenu.AddText("退出", nil, func(_ *menu.CallbackData) {
		runtime.Quit(app.ctx)
	})
	// }
	return appMenu
}

func runDesktop() {
	logrus.Infof("桌面模式启动中，端口: %d", port)
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             ProjectName,
		Width:             1024,
		Height:            768,
		Menu:              newMenu(app),
		HideWindowOnClose: true, // 关闭窗口时隐藏应用
		OnStartup:         app.startup,
		OnShutdown:        app.shutdown,
		AssetServer: &assetserver.Options{
			Assets:  web.WebDist,
			Handler: router.InitRouter(isDev, nil),
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			BackdropType:         windows.Mica, // 使用Mica效果
			Theme:                windows.SystemDefault,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			// About: &mac.AboutInfo{
			// 	Title:   ProjectName,
			// 	Message: "Copyright © 2025 AKimioJR(akimio.jr@gmail.com)",
			// 	Icon:    web.GetLogoSVGData(),
			// },
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
	if isServer { // 启动服务器模式
		runServer()
	} else { // 启动桌面模式
		runDesktop()
	}
}
