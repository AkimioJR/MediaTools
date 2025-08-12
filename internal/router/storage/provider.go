package storage

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /storage
// @Route /provider [get]
// @Summary 获取存储提供者列表
// @Description 返回所有已注册的存储提供者列表
// @Tags storage
// @Products json
// @Success 200 {object} schemas.Response[[]schemas.StorageProviderItem]
// @Failure 500 {object} schemas.Response[[]schemas.StorageProviderItem]
func ProviderList(ctx *gin.Context) {
	var resp schemas.Response[[]schemas.StorageProviderItem]

	resp.Data = storage_controller.ListStorageProviders()
	resp.Success = true
	logrus.Info(resp)
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage/provider
// @Route /{storage_type} [get]
// @Summary 获取指定存储提供者
// @Description 获取指定类型的存储提供者信息
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Accept json
// @Products json
// @Success 200 {object} schemas.Response[*schemas.StorageProviderItem]
// @Failure 400 {object} schemas.Response[*schemas.StorageProviderItem]
// @Failure 500 {object} schemas.Response[*schemas.StorageProviderItem]
func ProviderGet(ctx *gin.Context) {
	var resp schemas.Response[*schemas.StorageProviderItem]

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		logrus.Warningf("未知的存储类型: %s", storageTypeStr)
		resp.Message = "未知的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	item, err := storage_controller.GetStorageProvider(storageType)
	if err != nil {
		logrus.Warningf("获取存储提供者失败: %v", err)
		resp.Message = "获取存储提供者失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Success = true
	resp.Data = item
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage/provider
// @Route /{storage_type} [post]
// @Summary 注册新的存储器
// @Description 注册一个新的存储器
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param body body map[string]string true "存储器配置"
// @Accept json
// @Products json
// @Success 200 {object} schemas.Response[*schemas.StorageProviderItem]
// @Failure 400 {object} schemas.Response[*schemas.StorageProviderItem]
// @Failure 500 {object} schemas.Response[*schemas.StorageProviderItem]
func ProviderRegister(ctx *gin.Context) {
	var (
		req  map[string]string
		resp schemas.Response[*schemas.StorageProviderItem]
	)

	logrus.Debugf("请求体: %+v", ctx.Request.Body)

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		logrus.Warningf("未知的存储类型: %s", storageTypeStr)
		resp.Message = "未知的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("解析请求参数失败: %v", err)
		resp.Message = "解析请求参数失败: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	c := config.StorageConfig{
		Type: storageType,
		Data: req,
	}

	logrus.Debugf("注册存储器: %s, 配置: %+v", storageType, c)

	item, err := storage_controller.RegisterStorageProvider(c)
	if err != nil {
		logrus.Errorf("注册存储器失败: %v", err)
		resp.Message = "注册存储器失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	logrus.Debugf("存储器注册成功: %+v", item)

	resp.Success = true
	resp.Data = item
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage/provider
// @Route /{storage_type} [delete]
// @Summary 删除存储器
// @Description 删除指定类型的存储器
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Accept json
// @Products json
// @Success 200 {object} schemas.Response[*schemas.StorageProviderItem]
// @Failure 400 {object} schemas.Response[*schemas.StorageProviderItem]
// @Failure 500 {object} schemas.Response[*schemas.StorageProviderItem]
// @Router /storage/provider/{storage_type} [delete]
func ProviderDelete(ctx *gin.Context) {
	var resp schemas.Response[*schemas.StorageProviderItem]

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	switch storageType {
	case schemas.StorageUnknown:
		logrus.Warningf("未知的存储类型: %s", storageTypeStr)
		resp.Message = "未知的存储类型"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	case schemas.StorageLocal:
		logrus.Warning("无法删除本地存储器")
		resp.Message = "无法删除本地存储器"
		ctx.JSON(http.StatusBadRequest, resp)
	}

	item, err := storage_controller.UnRegisterStorageProvider(storageType)
	if err != nil {
		logrus.Errorf("删除存储器失败: %v", err)
		resp.Message = "删除存储器失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	logrus.Debugf("已删除存储器: %s", storageType)
	resp.Success = true
	resp.Data = item
	ctx.JSON(http.StatusOK, resp)
}
