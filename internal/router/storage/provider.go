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
	var resp schemas.Response[[]schemas.StorageProviderItem]
	resp.Data = storage_controller.ListStorageProviders()
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/provider/{storage_type} [get]
// @Summary 获取指定存储提供者
// @Description 获取指定类型的存储提供者信息
// @Tags 存储,存储器
// @Param storage_type path string true "存储类型"
// @Accept json
// @Products json
func ProviderGet(ctx *gin.Context) {
	var resp schemas.Response[*schemas.StorageProviderItem]

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "未知的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	item, err := storage_controller.GetStorageProvider(storageType)
	if err != nil {
		resp.Message = "获取存储提供者失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Data = item
	resp.RespondJSON(ctx, http.StatusOK)
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
		req  map[string]string
		resp schemas.Response[*schemas.StorageProviderItem]
	)

	logrus.Debugf("请求体: %+v", ctx.Request.Body)

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "未知的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "解析请求参数失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	c := config.StorageConfig{
		Type: storageType,
		Data: req,
	}

	logrus.Debugf("注册存储器: %s, 配置: %+v", storageType, c)

	item, err := storage_controller.RegisterStorageProvider(c)
	if err != nil {
		resp.Message = "注册存储器失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	logrus.Debugf("存储器注册成功: %+v", item)

	resp.Data = item
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/provider/{storage_type} [delete]
// @Summary 删除存储器
// @Description 删除指定类型的存储器
// @Tags 存储,存储器
// @Param storage_type path string true "存储类型"
// @Accept json
// @Products json
func ProviderDelete(ctx *gin.Context) {
	var resp schemas.Response[*schemas.StorageProviderItem]

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	switch storageType {
	case schemas.StorageUnknown:
		resp.Message = "未知的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return

	case schemas.StorageLocal:
		resp.Message = "无法删除本地存储器"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	item, err := storage_controller.UnRegisterStorageProvider(storageType)
	if err != nil {
		resp.Message = "删除存储器失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	logrus.Debugf("已删除存储器: %s", storageType)
	
	resp.Data = item
	resp.RespondJSON(ctx, http.StatusOK)
}
