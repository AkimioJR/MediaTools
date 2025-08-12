package storage

import (
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
