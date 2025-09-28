package app

import (
	"MediaTools/internal/router"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
)

func waitSysSign() <-chan os.Signal {
	sysSignCh := make(chan os.Signal, 1)
	signal.Notify(sysSignCh, syscall.SIGINT, syscall.SIGTERM)
	return sysSignCh
}

func runServer() {
	ginR := router.InitRouter(isDev)
	sysCh := waitSysSign()
	errCh := make(chan error, 1)
	go func() {
		err := ginR.Run(":" + strconv.Itoa(int(port)))
		if err != nil {
			errCh <- fmt.Errorf("启动服务器失败: %v", err)
		}
	}()
	select {
	case err := <-errCh:
		if err != nil {
			logrus.Errorf("应用程序运行中发生错误: %v", err)
		}
	case sig := <-sysCh:
		logrus.Infof("收到系统信号: %v, 退出应用程序", sig)
	}
}
