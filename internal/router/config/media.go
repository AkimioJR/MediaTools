package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/media_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /config/media
// @Router libraries [get]
// @Summary 获取媒体库配置
// @Description 获取媒体库配置
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response[[]config.LibraryConfig]
// @Failure 400 {object} schemas.Response[[]config.LibraryConfig]
// @Failure 500 {object} schemas.Response[[]config.LibraryConfig]
func MediaLibrary(ctx *gin.Context) {
	var resp schemas.Response[[]config.LibraryConfig]
	resp.Data = config.Media.Libraries
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /config/media
// @Router libraries [post]
// @Summary 更新媒体库配置
// @Description 更新媒体库配置
// @Tags config
// @Accept json
// @Produce json
// @Param config body []config.LibraryConfig true "媒体库配置"
// @Success 200 {object} schemas.Response[[]config.LibraryConfig]
// @Failure 400 {object} schemas.Response[[]config.LibraryConfig]
// @Failure 500 {object} schemas.Response[[]config.LibraryConfig]
func UpdateMediaLibrary(ctx *gin.Context) {
	var (
		req  []config.LibraryConfig
		resp schemas.Response[[]config.LibraryConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	logrus.Debugf("开始更新媒体库配置: %+v", req)

	oldConfig := config.Media.Libraries
	config.Media.Libraries = req
	err = media_controller.Init()
	if err != nil {
		logrus.Errorf("初始化 Media 控制器失败: %v", err)
		resp.Message = "初始化 Media 控制器失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		goto initErr
	}

	logrus.Debugf("Media 控制器初始化成功: %+v", config.Media.Libraries)

	err = config.WriteConfig()
	if err != nil {
		resp.Message = "更新配置失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = config.Media.Libraries
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
	return
initErr:
	config.Media.Libraries = oldConfig
	media_controller.Init()
}

// @BasePath /config/media
// @Router format [get]
// @Summary 获取媒体格式配置
// @Description 获取媒体格式配置
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response[config.FormatConfig]
// @Failure 400 {object} schemas.Response[config.FormatConfig]
// @Failure 500 {object} schemas.Response[config.FormatConfig]
func MediaFormat(ctx *gin.Context) {
	var resp schemas.Response[config.FormatConfig]
	resp.Data = config.Media.Format
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /config/media
// @Router format [post]
// @Summary 更新媒体格式配置
// @Description 更新媒体格式配置
// @Tags config
// @Accept json
// @Produce json
// @Param config body config.FormatConfig true "媒体格式配置"
// @Success 200 {object} schemas.Response[config.FormatConfig]
// @Failure 400 {object} schemas.Response[config.FormatConfig]
// @Failure 500 {object} schemas.Response[config.FormatConfig]
func UpdateMediaFormat(ctx *gin.Context) {
	var (
		req  config.FormatConfig
		resp schemas.Response[config.FormatConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	oldConfig := config.Media.Format
	config.Media.Format = req
	err = media_controller.InitFormatTemplates()
	if err != nil {
		resp.Message = "初始化格式模板失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		goto initErr
	}
	err = config.WriteConfig()
	if err != nil {
		resp.Message = "更新配置失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = config.Media.Format
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
	return
initErr:
	config.Media.Format = oldConfig
	media_controller.Init()
}

// @BasePath /config/media
// @Router custom_word [get]
// @Summary 获取自定义词配置
// @Description 获取自定义词配置
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response[config.CustomWordConfig]
// @Failure 400 {object} schemas.Response[config.CustomWordConfig]
// @Failure 500 {object} schemas.Response[config.CustomWordConfig]
func CustomWord(ctx *gin.Context) {
	var resp schemas.Response[config.CustomWordConfig]
	resp.Data = config.Media.CustomWord
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /config/media
// @Router custom_word [post]
// @Summary 更新自定义词配置
// @Description 更新自定义词配置
// @Tags config
// @Accept json
// @Produce json
// @Param config body config.CustomWordConfig true "自定义词配置"
// @Success 200 {object} schemas.Response[config.CustomWordConfig]
// @Failure 400 {object} schemas.Response[config.CustomWordConfig]
// @Failure 500 {object} schemas.Response[config.CustomWordConfig]
func UpdateCustomWord(ctx *gin.Context) {
	var (
		req  config.CustomWordConfig
		resp schemas.Response[config.CustomWordConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	oldConfig := config.Media.CustomWord
	config.Media.CustomWord = req
	err = media_controller.InitCustomWord()
	if err != nil {
		resp.Message = "初始化自定义词失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		goto initErr
	}

	err = config.WriteConfig()
	if err != nil {
		resp.Message = "更新配置失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = config.Media.CustomWord
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
	return
initErr:
	config.Media.CustomWord = oldConfig
	media_controller.Init()
}
