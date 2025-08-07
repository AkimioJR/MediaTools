package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /config
// @Router /tmdb [get]
// @Summary 获取 TMDB 配置
// @Description 获取 TMDB 配置
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response[config.TMDBConfig]
// @Failure 400 {object} schemas.Response[config.TMDBConfig]
// @Failure 500 {object} schemas.Response[config.TMDBConfig]
func TMDB(ctx *gin.Context) {
	var resp schemas.Response[config.TMDBConfig]
	resp.Data = config.TMDB
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /config
// @Router /tmdb [post]
// @Summary 更新 TMDB 配置
// @Description 更新 TMDB 配置
// @Tags config
// @Accept json
// @Produce json
// @Param config body config.TMDBConfig true "TMDB 配置"
// @Success 200 {object} schemas.Response[config.TMDBConfig]
// @Failure 400 {object} schemas.Response[config.TMDBConfig]
// @Failure 500 {object} schemas.Response[config.TMDBConfig]
func UpdateTMDB(ctx *gin.Context) {
	var (
		req  config.TMDBConfig
		resp schemas.Response[config.TMDBConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	config.TMDB = req
	err = config.WriteConfig()
	if err != nil {
		resp.Message = "更新配置失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	err = tmdb_controller.Init()
	if err != nil {
		resp.Message = "初始化 TMDB 控制器失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = config.TMDB
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}
