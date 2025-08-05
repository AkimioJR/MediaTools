package router

import (
	"MediaTools/internal/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRecentLogs(ctx *gin.Context) {
	var resp Response[[]string]
	resp.Data = logging.GetRecentLogs()
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}
