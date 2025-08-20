package history

import "github.com/gin-gonic/gin"

// 注册历史记录相关路由
func RegisterHistoryRouter(router *gin.Engine) {
	historyRouter := router.Group("/history") // 历史记录相关路由
	{
		mediaHistoryRouter := historyRouter.Group("/media") // 媒体历史记录相关路由
		{
			mediaHistoryRouter.GET("/transfer", QueryMediaTransferHistory) // 查询媒体转移历史记录
		}
	}
}
