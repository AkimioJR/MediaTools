package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/logging"
	"MediaTools/internal/schemas"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /config
// @Router /log [get]
// @Summary 获取日志配置
// @Description 获取日志配置
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response[config.LogConfig]
// @Failure 400 {object} schemas.Response[config.LogConfig]
// @Failure 500 {object} schemas.Response[config.LogConfig]
func Log(ctx *gin.Context) {
	var resp schemas.Response[config.LogConfig]
	resp.Data = config.Log
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /config
// @Router /log [post]
// @Summary 更新日志配置
// @Description 更新日志配置
// @Tags config
// @Accept json
// @Produce json
// @Param config body config.LogConfig true "日志配置"
// @Success 200 {object} schemas.Response[config.LogConfig]
// @Failure 400 {object} schemas.Response[config.LogConfig]
// @Failure 500 {object} schemas.Response[config.LogConfig]
func UpdateLog(ctx *gin.Context) {
	var (
		req       config.LogConfig
		resp      schemas.Response[config.LogConfig]
		oldConfig = config.Log
		err       error
	)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	logrus.Debugf("开始更新日志配置: %+v", req)

	if config.Log.ConsoleLevel != req.ConsoleLevel {
		config.Log.ConsoleLevel = req.ConsoleLevel
		err = logging.SetLevel(req.ConsoleLevel)
		if err != nil {
			err = fmt.Errorf("设置终端日志级别失败: %w", err)
			goto responseErr
		}
	}
	if config.Log.FileLevel != req.FileLevel {
		config.Log.FileLevel = req.FileLevel
		err = logging.SetFileLevel(req.FileLevel)
		if err != nil {
			err = fmt.Errorf("设置文件日志级别失败: %w", err)
			goto responseErr
		}
	}
	if config.Log.FileDir != req.FileDir {
		config.Log.FileDir = req.FileDir
		logging.SetLogDir(req.FileDir) // 更新日志目录
	}

	logrus.Debugf("初始化日志配置成功: %+v", config.Log)

	logrus.Debug("开始更新配置文件")
	err = config.WriteConfig()
	if err != nil {
		resp.Message = "更新配置失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = config.Log
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
	return

responseErr:
	resp.Message = err.Error()
	logrus.Errorf("更新日志配置失败: %v", err)
	config.Log = oldConfig // 恢复旧配置
	logging.Init()         // 重新初始化日志系统
	logrus.Debugf("日志配置恢复成功: %+v", config.Log)
	ctx.JSON(http.StatusInternalServerError, resp)
}
