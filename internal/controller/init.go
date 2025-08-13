package controller

import (
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/controller/library_controller"
	"MediaTools/internal/controller/recognize_controller"
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"fmt"
)

type InitFunc func() error

var initFuncs = []InitFunc{
	tmdb_controller.Init,
	fanart_controller.Init,
	scrape_controller.Init,
	storage_controller.Init,
	library_controller.Init,
	recognize_controller.Init,
}

func InitAllControllers() error {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			return fmt.Errorf("初始化工具链失败: %w", err)
		}
	}
	return nil
}
