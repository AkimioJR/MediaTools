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
// @Tags config
// @Produce json
// @Success 200 {object} config.FanartConfig
func Fanart(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, config.Fanart)
}

// @Router /config/fanart [post]
// @Summary 更新 Fanart 配置
// @Description 更新 Fanart 配置
// @Tags config
// @Accept json
// @Produce json
// @Param config body config.FanartConfig true "Fanart 配置"
// @Success 200 {object} config.FanartConfig
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func UpdateFanart(ctx *gin.Context) {
	var req config.FanartConfig

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, schemas.ErrResponse{Message: "请求参数错误: " + err.Error()})
		return
	}

	logrus.Debugf("开始更新 Fanart 配置: %+v", req)

	oldConfig := config.Fanart
	config.Fanart = req
	err = fanart_controller.Init()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, schemas.ErrResponse{Message: "初始化 Fanart 控制器失败: " + err.Error()})
		goto initErr
	}

	logrus.Debugf("Fanart 控制器初始化成功: %+v", config.Fanart)

	err = config.WriteConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schemas.ErrResponse{Message: "写入配置文件失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, config.Fanart)
	return

initErr:
	config.Fanart = oldConfig
	fanart_controller.Init()
}
