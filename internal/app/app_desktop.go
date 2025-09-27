//go:build !onlyServer

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"
	webview "github.com/webview/webview_go"
)

const AppName = "MediaTools"

var SupportDesktopMode = true

func openWindows(port uint) <-chan struct{} {
	doneCh := make(chan struct{})
	go func() {
		w := webview.New(false)
		defer w.Destroy()
		w.SetSize(800, 600, webview.HintNone)
		w.SetTitle(AppName)
		w.Navigate(fmt.Sprintf("http://localhost:%d", port))
		w.Run()
		close(doneCh)
	}()
	return doneCh
}

func runDesktop() <-chan error {
	logrus.Infof("桌面模式启动中，端口: %d", port)
	serverChan := runServer() // 在后台启动服务器
	windowsCh := openWindows(port)
	errCh := make(chan error)
	go func() {
		select {
		case err := <-serverChan:
			errCh <- fmt.Errorf("服务器运行中发生错误: %v", err)
		case <-windowsCh:
			errCh <- nil
		}
		close(errCh)
	}()
	return errCh
}

func Run() <-chan error {
	if isServer { // 启动服务器模式
		return runServer()
	} else { // 启动桌面模式
		return runDesktop()
	}
}
