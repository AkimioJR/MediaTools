package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Router /config/fanart [get]
// @Summary 获取 Fanart 配置
// @Description 获取 Fanart 配置
// @Tags 应用配置,Fanart
// @Produce json
func Fanart(ctx *gin.Context) {
	var resp schemas.Response[*config.FanartConfig]
	resp.Success = true
	resp.Data = &config.Fanart
	logrus.Debugf("获取 Fanart 配置: %+v", resp.Data)
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Router /config/fanart [post]
// @Summary 更新 Fanart 配置
// @Description 更新 Fanart 配置
// @Tags 应用配置,Fanart
// @Accept json
// @Produce json
// @Param config body config.FanartConfig true "Fanart 配置"
func UpdateFanart(ctx *gin.Context) {
	var (
		req  config.FanartConfig
		resp schemas.Response[*config.FanartConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		logrus.Warning(resp.Message)
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	logrus.Debugf("开始更新 Fanart 配置: %+v", req)

	oldConfig := config.Fanart
	config.Fanart = req
	err = fanart_controller.Init()
	if err != nil {
		resp.Message = "初始化 Fanart 控制器失败: " + err.Error()
		logrus.Warning(resp.Message)
		logrus.Debugf("开始恢复 Fanart 配置: %+v", oldConfig)
		config.Fanart = oldConfig
		fanart_controller.Init()
		logrus.Debugf("恢复 Fanart 配置成功: %+v", config.Fanart)
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	logrus.Debugf("Fanart 控制器初始化成功: %+v", config.Fanart)

	err = config.WriteConfig()
	if err != nil {
		resp.Message = "写入配置文件失败: " + err.Error()
		logrus.Warning(resp.Message)
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Success = true
	resp.Data = &config.Fanart
	ctx.JSON(http.StatusOK, resp)
}
