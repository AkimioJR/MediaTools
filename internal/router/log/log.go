package log

import (
	"MediaTools/internal/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /log
// @Router /recent [get]
// @Summary 获取最近日志
// @Description 获取最近日志
// @Tags log
// @Produce json
// @Success 200 {object} []loghook.LogDetail
func GetRecentLogs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, logging.GetRecentLogs())
}
