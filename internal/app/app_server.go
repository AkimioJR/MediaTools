//go:build onlyServer

package app

import "github.com/sirupsen/logrus"

var (
	SupportDesktopMode = false
)

func Run() {
	serverCh := runServer()
	
	select {
	case err := <-serverCh:
		if err != nil {
			logrus.Errorf("应用程序运行中发生错误: %v", err)
		}
	case sig := <-sysCh:
		logrus.Infof("收到系统信号: %v, 退出应用程序", sig)
	}
}
