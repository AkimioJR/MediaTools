package library

import "github.com/gin-gonic/gin"

func RegisterLibraryRouter(router *gin.Engine) {
	libraryRouter := router.Group("/library") // 媒体库相关接口
	{
		libraryRouter.POST("/archive", LibraryArchiveMedia) // 归档媒体文件
	}
}
