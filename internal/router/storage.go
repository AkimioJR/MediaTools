package router

import (
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StorageList(c *gin.Context) {
	var resp Response[[]schemas.StorageProviderItem]
	resp.Data = storage_controller.ListStorageProviders()
	resp.Success = true
	logrus.Info(resp)
	c.JSON(http.StatusOK, resp)
}
