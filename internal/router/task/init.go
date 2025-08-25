package task

import "github.com/gin-gonic/gin"

func RegisterTaskRouter(taskRouter *gin.RouterGroup) {

	transferRouter := taskRouter.Group("/transfer") // 媒体转移任务相关接口
	{
		transferRouter.GET("/", GetAllTransferTasks)      // 查询转移任务列表
		transferRouter.GET("/:id", GetTransferTask)       // 获取转移任务状态
		transferRouter.DELETE("/:id", CancelTransferTask) // 取消转移任务
	}

}
