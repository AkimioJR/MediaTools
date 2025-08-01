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
		storageRouter.GET("/list", StorageProviderList) // 基础信息接口

		storageTypeRouter := storageRouter.Group("/:storage_type") // 按存储类型分组的API
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

		storageRouter.POST("/copy", StorageCopyFile)         // 复制文件
		storageRouter.POST("/move", StorageMoveFile)         // 移动文件
		storageRouter.POST("/link", StorageLinkFile)         // 创建硬链接
		storageRouter.POST("/softlink", StorageSoftLinkFile) // 创建软链接
		storageRouter.POST("/transfer", StorageTransferFile) // 通用文件传输接口
	}

	return ginRouter
}
