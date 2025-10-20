//go:build desktop && !darwin
// +build desktop,!darwin

package app

import (
	"MediaTools/internal/info"
	"MediaTools/web"

	"fyne.io/systray"
)

func (a *App) onReadyFunc() {
	systray.SetIcon(web.GetIconData())

	mAbout := systray.AddMenuItem("关于", "关于 "+info.ProjectName)

	systray.AddSeparator()

	mShowWindow := systray.AddMenuItem(ShowWindowsString, ShowWindowsTip)
	mHideWindow := systray.AddMenuItem(HideWindowsString, HideWindowsTip)

	systray.AddSeparator()

	mQuit := systray.AddMenuItem(QuitString, QuitTip)
	// mswith := systray.AddMenuItem("switch", "switch")

	go func() {
		for {
			select {
			case <-mAbout.ClickedCh:
				a.OpenAboutDialog()

			case <-mShowWindow.ClickedCh:
				a.Show()
			case <-mHideWindow.ClickedCh:
				a.Hide()

			case <-mQuit.ClickedCh:
				a.Quit()

				// case <-mswith.ClickedCh:
				// 	if runtime.WindowIsFullscreen(a.ctx) {
				// 		runtime.WindowUnfullscreen(a.ctx)
				// 	} else {
				// 		runtime.WindowFullscreen(a.ctx)
				// 	}
			}
		}
	}()
}

func (a *App) createSystray() {
	startFunc, endFunc := systray.RunWithExternalLoop(a.onReadyFunc, nil)
	a.systrayEndfunc = endFunc
	startFunc()
}
