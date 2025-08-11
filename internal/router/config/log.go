package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/logging"
	"MediaTools/internal/schemas"
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
		req  config.LogConfig
		resp schemas.Response[config.LogConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	oldConfig := config.Log
	config.Log = req
	logrus.Debug("开始更新日志配置: ", req)
	err = logging.Init()
	if err != nil {
		logrus.Error("更新日志配置失败: ", err)
		resp.Message = "更新日志配置失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		config.Log = oldConfig
		logging.Init()
	}
	logrus.Debug("日志配置更新成功: ", config.Log)

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
}
