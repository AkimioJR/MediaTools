package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /config
// @Router /fanart [get]
// @Summary 获取 Fanart 配置
// @Description 获取 Fanart 配置
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response[config.FanartConfig]
// @Failure 400 {object} schemas.Response[config.FanartConfig]
// @Failure 500 {object} schemas.Response[config.FanartConfig]
func Fanart(ctx *gin.Context) {
	var resp schemas.Response[config.FanartConfig]
	resp.Data = config.Fanart
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /config
// @Router /fanart [post]
// @Summary 更新 Fanart 配置
// @Description 更新 Fanart 配置
// @Tags config
// @Accept json
// @Produce json
// @Param config body config.FanartConfig true "Fanart 配置"
// @Success 200 {object} schemas.Response[config.FanartConfig]
// @Failure 400 {object} schemas.Response[config.FanartConfig]
// @Failure 500 {object} schemas.Response[config.FanartConfig]
func UpdateFanart(ctx *gin.Context) {
	var (
		req  config.FanartConfig
		resp schemas.Response[config.FanartConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	oldConfig := config.Fanart
	config.Fanart = req
	err = fanart_controller.Init()
	if err != nil {
		resp.Message = "初始化 Fanart 控制器失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		goto initErr
	}

	err = config.WriteConfig()
	if err != nil {
		resp.Message = "更新配置失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = config.Fanart
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
	return
initErr:
	config.Fanart = oldConfig
	fanart_controller.Init()
}
