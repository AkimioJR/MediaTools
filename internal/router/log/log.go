package log

import (
	"MediaTools/internal/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Router /log/recent [get]
// @Summary 获取最近日志
// @Description 获取最近日志
// @Tags 日志
// @Produce json
func GetRecentLogs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, logging.GetRecentLogs())
}
