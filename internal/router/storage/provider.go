package storage

import (
	"MediaTools/internal/config"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Route /storage/provider [get]
// @Summary 获取存储提供者列表
// @Description 返回所有已注册的存储提供者列表
// @Tags 存储,存储器
// @Products json
func ProviderList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, storage_controller.ListStorageProviders())
}

// @Route /storage/provider/{storage_type} [get]
// @Summary 获取指定存储提供者
// @Description 获取指定类型的存储提供者信息
// @Tags 存储,存储器
// @Param storage_type path string true "存储类型"
// @Accept json
// @Products json
func ProviderGet(ctx *gin.Context) {
	var errResp schemas.ErrResponse

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "未知的存储类型: " + storageTypeStr
		logrus.Warning(errResp.Message)
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	item, err := storage_controller.GetStorageProvider(storageType)
	if err != nil {
		errResp.Message = "获取存储提供者失败: " + err.Error()
		logrus.Warning(errResp.Message)
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// @Route /storage/provider/{storage_type} [post]
// @Summary 注册新的存储器
// @Description 注册一个新的存储器
// @Tags 存储,存储器
// @Param storage_type path string true "存储类型"
// @Param body body map[string]string true "存储器配置"
// @Accept json
// @Products json
func ProviderRegister(ctx *gin.Context) {
	var (
		req     map[string]string
		errResp schemas.ErrResponse
	)

	logrus.Debugf("请求体: %+v", ctx.Request.Body)

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "未知的存储类型: " + storageTypeStr
		logrus.Warning(errResp.Message)
		ctx.JSON(http.StatusBadRequest, errResp)
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errResp.Message = "解析请求参数失败: " + err.Error()
		logrus.Warning(errResp.Message)
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	c := config.StorageConfig{
		Type: storageType,
		Data: req,
	}

	logrus.Debugf("注册存储器: %s, 配置: %+v", storageType, c)

	item, err := storage_controller.RegisterStorageProvider(c)
	if err != nil {
		errResp.Message = "注册存储器失败: " + err.Error()
		logrus.Warning(errResp.Message)
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	logrus.Debugf("存储器注册成功: %+v", item)
	ctx.JSON(http.StatusOK, item)
}

// @Route /storage/provider/{storage_type} [delete]
// @Summary 删除存储器
// @Description 删除指定类型的存储器
// @Tags 存储,存储器
// @Param storage_type path string true "存储类型"
// @Accept json
// @Products json
func ProviderDelete(ctx *gin.Context) {
	var errResp schemas.ErrResponse

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	switch storageType {
	case schemas.StorageUnknown:
		errResp.Message = "未知的存储类型"
		logrus.Warning(errResp.Message)
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	case schemas.StorageLocal:
		errResp.Message = "无法删除本地存储器"
		logrus.Warning(errResp.Message)
		ctx.JSON(http.StatusBadRequest, errResp)
	}

	item, err := storage_controller.UnRegisterStorageProvider(storageType)
	if err != nil {
		errResp.Message = "删除存储器失败: " + err.Error()
		logrus.Warning(errResp.Message)
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	logrus.Debugf("已删除存储器: %s", storageType)
	ctx.JSON(http.StatusOK, item)
}
