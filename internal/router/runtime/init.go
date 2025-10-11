package runtime

import (
	"MediaTools/internal/info"
	"MediaTools/internal/schemas"

	"github.com/gin-gonic/gin"
)

func RegisterRuntimeRouter(r *gin.RouterGroup) {
	r.GET("/status", Status) // 获取程序运行状态
}

// @Router /runtime/status [get]
// @Summary 获取程序运行状态
// @Description 获取程序运行状态信息
// @Tags App
// @Produce json
func Status(ctx *gin.Context) {
	var resp schemas.Response[*info.RuntimeAppStatusInfo]
	resp.RespondSuccessJSON(ctx, &info.RuntimeAppStatus)
}
