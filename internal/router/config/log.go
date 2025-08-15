package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/logging"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Router /config/log [get]
// @Summary 获取日志配置
// @Description 获取日志配置
// @Tags 应用配置,日志
// @Produce json
func Log(ctx *gin.Context) {
	var resp schemas.Response[*config.LogConfig]
	resp.RespondSuccessJSON(ctx, &config.Log)
}

// @Router /config/log [post]
// @Summary 更新日志配置
// @Description 更新日志配置
// @Tags 应用配置,日志
// @Accept json
// @Produce json
// @Param config body config.LogConfig true "日志配置"
func UpdateLog(ctx *gin.Context) {
	var (
		req       config.LogConfig
		oldConfig = config.Log
		resp      schemas.Response[*config.LogConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	logrus.Debugf("开始更新日志配置: %+v", req)

	if config.Log.ConsoleLevel != req.ConsoleLevel {
		config.Log.ConsoleLevel = req.ConsoleLevel
		err = logging.SetLevel(req.ConsoleLevel)
		if err != nil {
			resp.Message = "设置日志级别失败: " + err.Error()
			resp.RespondJSON(ctx, http.StatusInternalServerError)
			logging.SetLevel(oldConfig.ConsoleLevel) // 恢复旧级别
			return
		}
	}
	if config.Log.FileLevel != req.FileLevel {
		config.Log.FileLevel = req.FileLevel
		err = logging.SetFileLevel(req.FileLevel)
		if err != nil {
			resp.Message = "设置文件日志级别失败: " + err.Error()
			resp.RespondJSON(ctx, http.StatusInternalServerError)
			logging.SetFileLevel(oldConfig.FileLevel) // 恢复旧级别
			logrus.Debugf("恢复旧的文件日志级别: %s", oldConfig.FileLevel)
			return
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
		resp.Message = "写入新配置文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Success = true
	resp.Data = &config.Log
	resp.RespondJSON(ctx, http.StatusOK)
}
