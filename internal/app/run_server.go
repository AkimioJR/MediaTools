package app

import (
	"MediaTools/internal/router"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

func runServer() <-chan error {
	ginR := router.InitRouter(isDev, webDist)
	errCh := make(chan error, 1)
	// 在服务器构建中总是启动服务器模式
	go func() {
		err := ginR.Run(":" + strconv.Itoa(int(port)))
		if err != nil {
			errCh <- fmt.Errorf("启动服务器失败: %v", err)
		}
	}()
	logrus.Infof("服务器启动成功，监听端口: %d", port)
	return errCh
}
