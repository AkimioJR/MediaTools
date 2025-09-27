package app

import (
	"MediaTools/internal/router"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

func runServer() {
	ginR := router.InitRouter(isDev, webDist)

	// 在服务器构建中总是启动服务器模式
	err := ginR.Run(":" + strconv.Itoa(int(port)))
	if err != nil {
		panic(fmt.Sprintf("启动服务器失败: %v", err))
	}
	logrus.Infof("服务器启动成功，监听端口: %d", port)
}
