package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Router /config/tmdb [get]
// @Summary 获取 TMDB 配置
// @Description 获取 TMDB 配置
// @Tags 应用配置,TMDB
// @Produce json
// @Success 200 {object} config.TMDBConfig
func TMDB(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, config.TMDB)
}

// @Router /config/tmdb [post]
// @Summary 更新 TMDB 配置
// @Description 更新 TMDB 配置
// @Tags 应用配置,TMDB
// @Accept json
// @Param config body config.TMDBConfig true "TMDB 配置"
// @Success 200 {object} config.TMDBConfig
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func UpdateTMDB(ctx *gin.Context) {
	var (
		req     config.TMDBConfig
		errResp schemas.ErrResponse
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errResp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	logrus.Debugf("开始更新 TMDB 配置: %+v", req)
	oldConfig := config.TMDB
	config.TMDB = req

	err = tmdb_controller.Init()
	if err != nil {
		logrus.Errorf("初始化 TMDB 控制器失败: %v", err)
		errResp.Message = "初始化 TMDB 控制器失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		goto initErr
	}

	logrus.Debugf("TMDB 控制器初始化成功: %+v", config.TMDB)

	err = config.WriteConfig()
	if err != nil {
		errResp.Message = "更新配置失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, config.TMDB)
	return

initErr:
	config.TMDB = oldConfig
	tmdb_controller.Init()
}
