package history

import "github.com/gin-gonic/gin"

// 注册历史记录相关路由
func RegisterHistoryRouter(historyRouter *gin.RouterGroup) {
	mediaHistoryRouter := historyRouter.Group("/media") // 媒体历史记录相关路由
	{
		mediaHistoryRouter.GET("", QueryMediaTransferHistory)         // 查询媒体转移历史记录
		mediaHistoryRouter.GET("/:id", QueryMediaTransferHistoryByID) // 查询媒体转移历史记录 by ID
		mediaHistoryRouter.DELETE("/:id", DeleteMediaTransferHistory) // 删除媒体转移历史记录
	}
}
