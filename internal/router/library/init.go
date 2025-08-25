package library

import "github.com/gin-gonic/gin"

func RegisterLibraryRouter(libraryRouter *gin.RouterGroup) {
	libraryRouter.POST("/archive", ArchiveMediaManual) // 手动归档媒体文件
}
