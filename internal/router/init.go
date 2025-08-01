package router

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	ginRouter := gin.Default()
	mediaRouter := ginRouter.Group("/media")
	{
		mediaRouter.GET("/recognize", MediaRecognize)
	}
	storageRouter := ginRouter.Group("/storage")
	{
		storageRouter.GET("/list", StorageList)
	}

	return ginRouter
}
