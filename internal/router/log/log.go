package log

import (
	"MediaTools/internal/logging"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /log
// @Router /recent [get]
// @Summary 获取最近日志
// @Description 获取最近日志
// @Tags log
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response[[]string]
// @Failure 400 {object} schemas.Response[[]string]
// @Failure 500 {object} schemas.Response[[]string]
func GetRecentLogs(ctx *gin.Context) {
	var resp schemas.Response[[]string]
	resp.Data = logging.GetRecentLogs()
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}
