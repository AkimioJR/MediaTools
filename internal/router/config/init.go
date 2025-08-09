package config

import (
	"github.com/gin-gonic/gin"
)

// 注册配置相关路由
func RegisterConfigRouter(router *gin.Engine) {
	configRouter := router.Group("/config")
	{
		configRouter.GET("/log", Log)
		configRouter.POST("/log", UpdateLog)

		configRouter.GET("/tmdb", TMDB)
		configRouter.POST("/tmdb", UpdateTMDB)

		configRouter.GET("/fanart", Fanart)
		configRouter.POST("/fanart", UpdateFanart)

		mediaRouter := configRouter.Group("/media")
		{
			mediaRouter.GET("/libraries", MediaLibrary)
			mediaRouter.POST("/libraries", UpdateMediaLibrary)

			mediaRouter.GET("/format", MediaFormat)
			mediaRouter.POST("/format", UpdateMediaFormat)

			mediaRouter.GET("/custom_word", CustomWord)
			mediaRouter.POST("/custom_word", UpdateCustomWord)
		}
	}
}
