package log

import "github.com/gin-gonic/gin"

// 注册日志相关路由
func RegisterLogRouter(logRouter *gin.RouterGroup) {
	logRouter.GET("/recent", GetRecentLogs) // 获取最近日志
}
