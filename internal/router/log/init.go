package log

import "github.com/gin-gonic/gin"

// 注册日志相关路由
func RegisterLogRouter(router *gin.Engine) {
	logRouter := router.Group("/log") // 日志相关接口
	{
		logRouter.GET("/recent", GetRecentLogs) // 获取最近日志
	}
}
