package controller

import (
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/controller/transfer_controller"
	"fmt"

	"github.com/sirupsen/logrus"
)

type InitFunc func() error

var initFuncs = []InitFunc{
	tmdb_controller.Init,
	fanart_controller.Init,
	scrape_controller.Init,
	transfer_controller.Init,
}

func InitAllControllers() error {
	logrus.Info("开始初始化全部工具链...")
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			return fmt.Errorf("初始化工具链失败: %w", err)
		}
	}
	logrus.Info("全部工具链初始化完成")
	return nil
}
