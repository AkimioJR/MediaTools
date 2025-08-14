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
func TMDB(ctx *gin.Context) {
	var resp schemas.Response[*config.TMDBConfig]
	resp.Success = true
	resp.Data = &config.TMDB
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Router /config/tmdb [post]
// @Summary 更新 TMDB 配置
// @Description 更新 TMDB 配置
// @Tags 应用配置,TMDB
// @Accept json
// @Param config body config.TMDBConfig true "TMDB 配置"
func UpdateTMDB(ctx *gin.Context) {
	var (
		req  config.TMDBConfig
		resp schemas.Response[*config.TMDBConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	logrus.Debugf("开始更新 TMDB 配置: %+v", req)
	oldConfig := config.TMDB
	config.TMDB = req

	err = tmdb_controller.Init()
	if err != nil {
		logrus.Errorf("初始化 TMDB 控制器失败: %v", err)
		resp.Message = "初始化 TMDB 控制器失败: " + err.Error()
		config.TMDB = oldConfig
		tmdb_controller.Init()
		logrus.Debugf("恢复 TMDB 原始数据: %+v", config.TMDB)
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	logrus.Debugf("TMDB 控制器初始化成功: %+v", config.TMDB)

	err = config.WriteConfig()
	if err != nil {
		resp.Message = "写入新配置文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.RespondJSON(ctx, http.StatusOK)
}
