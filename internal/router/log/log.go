package log

import (
	"MediaTools/internal/logging"
	"MediaTools/internal/pkg/loghook"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Router /log/recent [get]
// @Summary 获取最近日志
// @Description 获取最近日志
// @Tags 日志
// @Produce json
func GetRecentLogs(ctx *gin.Context) {
	var resp schemas.Response[[]loghook.LogDetail]
	resp.Success = true
	resp.Data = logging.GetRecentLogs()
	resp.RespondJSON(ctx, http.StatusOK)
}
