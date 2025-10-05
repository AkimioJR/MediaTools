package server

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller"
	"MediaTools/internal/database"
	"MediaTools/internal/logging"
	"fmt"

	"github.com/sirupsen/logrus"
)

func InitApp(isDev bool, port uint16) {
	logrus.Info("开始初始化配置...")
	err := config.Init()
	if err != nil {
		panic(fmt.Sprintf("配置初始化失败: %v", err))
	}
	logrus.Info("配置初始化完成")

	logrus.Info("开始初始化日志...")
	err = logging.Init()
	if err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
	logrus.Info("日志初始化完成")
	err = database.Init()
	if err != nil {
		panic(fmt.Sprintf("数据库初始化失败: %v", err))
	}

	logrus.Info("开始初始化全部工具链...")
	err = controller.InitAllControllers()
	if err != nil {
		panic(fmt.Sprintf("工具链初始化失败: %v", err))
	}
	logrus.Info("全部工具链初始化完成")

}
