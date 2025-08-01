package router

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	ginRouter := gin.Default()

	mediaRouter := ginRouter.Group("/media") // 媒体相关接口
	{
		mediaRouter.GET("/recognize", MediaRecognize)
	}

	scrapeRouter := ginRouter.Group("/scrape") // 刮削相关接口
	{
		scrapeRouter.POST("/video", ScrapeVideo) // 刮削视频
	}

	storageRouter := ginRouter.Group("/storage") // 存储相关接口
	{
		// 基础信息接口
		storageRouter.GET("/list", StorageProviderList)

		// 按存储类型分组的API
		storageTypeRouter := storageRouter.Group("/:storage_type")
		{
			// 基础操作接口
			storageTypeRouter.GET("/info", StorageGetFileInfo)
			storageTypeRouter.GET("/exists", StorageCheckExists)
			storageTypeRouter.GET("/list", StorageList)

			// 文件和目录操作接口
			storageTypeRouter.POST("/mkdir", StorageMkdir)
			storageTypeRouter.DELETE("/delete", StorageDelete)

			// 文件传输接口
			storageTypeRouter.POST("/upload", StorageUploadFile)
			storageTypeRouter.GET("/download", StorageDownloadFile)
		}

		storageRouter.POST("/copy", StorageCopyFile)
		storageRouter.POST("/move", StorageMoveFile)
		storageRouter.POST("/link", StorageLinkFile)
		storageRouter.POST("/softlink", StorageSoftLinkFile)
		storageRouter.POST("/transfer", StorageTransferFile)
	}

	return ginRouter
}
