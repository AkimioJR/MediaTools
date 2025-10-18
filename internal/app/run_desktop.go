//go:build desktop
// +build desktop

package app

import (
	"MediaTools/internal/info"
	"MediaTools/internal/router"
	"MediaTools/web"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func init() {
	info.Version.SupportDesktopMode = true
}

func runDesktop() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             info.ProjectName,
		Width:             1024,
		Height:            768,
		Menu:              app.newMenu(),
		Frameless:         info.Version.OS == "windows", // Windows下使用无边框窗口
		HideWindowOnClose: true,                         // 关闭窗口时隐藏应用
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
