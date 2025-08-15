package config

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/recognize_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Router /config/libraries [get]
// @Summary 获取媒体库配置
// @Description 获取媒体库配置
// @Tags 应用配置
// @Produce json
func MediaLibrary(ctx *gin.Context) {
	var resp schemas.Response[[]config.LibraryConfig]
	resp.Success = true
	resp.Data = config.Media.Libraries
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Router /config/libraries [post]
// @Summary 更新媒体库配置
// @Description 更新媒体库配置
// @Tags 应用配置
// @Accept json
// @Produce json
// @Param config body []config.LibraryConfig true "媒体库配置"
func UpdateMediaLibrary(ctx *gin.Context) {
	var (
		req  []config.LibraryConfig
		resp schemas.Response[[]config.LibraryConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	logrus.Debugf("开始更新媒体库配置: %+v", req)

	oldConfig := config.Media.Libraries
	config.Media.Libraries = req
	err = recognize_controller.Init()
	if err != nil {
		logrus.Warningf("初始化 Media 控制器失败: %v", err)
		resp.Message = "初始化 Media 控制器失败: " + err.Error()
		config.Media.Libraries = oldConfig
		recognize_controller.Init()
		logrus.Debugf("恢复旧的媒体库配置: %+v", config.Media.Libraries)
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	logrus.Debugf("Media 控制器初始化成功: %+v", config.Media.Libraries)

	err = config.WriteConfig()
	if err != nil {
		resp.Message = "写入配置文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Success = true
	resp.Data = config.Media.Libraries
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Router /config/media/format [get]
// @Summary 获取媒体格式配置
// @Description 获取媒体格式配置
// @Tags 应用配置
// @Produce json
func MediaFormat(ctx *gin.Context) {
	var resp schemas.Response[*config.FormatConfig]
	resp.Success = true
	resp.Data = &config.Media.Format
	logrus.Debugf("获取媒体格式配置: %+v", resp.Data)
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Router /config/media/format [post]
// @Summary 更新媒体格式配置
// @Description 更新媒体格式配置
// @Tags 应用配置
// @Accept json
// @Produce json
// @Param config body config.FormatConfig true "媒体格式配置"
func UpdateMediaFormat(ctx *gin.Context) {
	var (
		req  config.FormatConfig
		resp schemas.Response[*config.FormatConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		logrus.Warning(resp.Message)
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	oldConfig := config.Media.Format
	config.Media.Format = req
	err = recognize_controller.InitFormatTemplates()
	if err != nil {
		resp.Message = "初始化格式模板失败: " + err.Error()
		config.Media.Format = oldConfig
		recognize_controller.Init()
		logrus.Debugf("恢复旧的媒体格式配置: %+v", config.Media.Format)
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	err = config.WriteConfig()
	if err != nil {
		resp.Message = "写入配置文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Success = true
	resp.Data = &config.Media.Format
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Router /config/media/custom_word [get]
// @Summary 获取自定义词配置
// @Description 获取自定义词配置
// @Tags 应用配置
// @Accept json
// @Produce json
func CustomWord(ctx *gin.Context) {
	var resp schemas.Response[*config.CustomWordConfig]
	resp.Success = true
	resp.Data = &config.Media.CustomWord
	logrus.Debugf("获取自定义词配置: %+v", resp.Data)
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Router /config/media/custom_word [post]
// @Summary 更新自定义词配置
// @Description 更新自定义词配置
// @Tags 应用配置
// @Accept json
// @Produce json
// @Param config body config.CustomWordConfig true "自定义词配置"
func UpdateCustomWord(ctx *gin.Context) {
	var (
		req  config.CustomWordConfig
		resp schemas.Response[*config.CustomWordConfig]
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		logrus.Warning(resp.Message)
		resp.RespondJSON(ctx, http.StatusBadGateway)
		return
	}

	oldConfig := config.Media.CustomWord
	config.Media.CustomWord = req
	err = recognize_controller.InitCustomWord()
	if err != nil {
		resp.Message = "初始化自定义词失败: " + err.Error()
		logrus.Warning(resp.Message)
		config.Media.CustomWord = oldConfig
		recognize_controller.Init()
		logrus.Debugf("恢复旧的自定义词配置: %+v", config.Media.CustomWord)
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	err = config.WriteConfig()
	if err != nil {
		resp.Message = "写入配置文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Success = true
	resp.Data = &config.Media.CustomWord
	resp.RespondJSON(ctx, http.StatusOK)
}
