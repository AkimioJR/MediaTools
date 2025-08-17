package library

import "github.com/gin-gonic/gin"

func RegisterLibraryRouter(router *gin.Engine) {
	libraryRouter := router.Group("/library") // 媒体库相关接口
	{
		libraryRouter.POST("/archive", ArchiveMediaManual) // 手动归档媒体文件
	}
}
