package storage

import "github.com/gin-gonic/gin"

func RegisterStorageRouter(router *gin.Engine) {
	storageRouter := router.Group("/storage") // 存储相关接口
	{
		providerRouter := storageRouter.Group("/provider") // 存储提供者相关接口
		{
			providerRouter.GET("", ProviderList)                    // 获取存储提供者列表
			providerRouter.GET("/:storage_type", ProviderGet)       // 获取指定存储提供者
			providerRouter.POST("/:storage_type", ProviderRegister) // 注册新的存储提供者
			providerRouter.DELETE("/:storage_type", ProviderDelete) // 删除存储提供者
		}

		storageTypeRouter := storageRouter.Group("/:storage_type") // 按存储类型分组的API
		{
			// 基础操作接口
			storageTypeRouter.GET("/info", StorageGetFileInfo)
			storageTypeRouter.GET("/exists", StorageCheckExists)
			storageTypeRouter.GET("/list", StorageList)              // 列出目录内容（非详细信息）
			storageTypeRouter.GET("/list_detail", StorageListDetail) // 列出目录内容（详细信息）

			// 文件和目录操作接口
			storageTypeRouter.POST("/mkdir", StorageMkdir)
			storageTypeRouter.POST("/rename", StorageRename)
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

}
