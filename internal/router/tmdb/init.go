package tmdb

import "github.com/gin-gonic/gin"

// TMDB 相关接口
func RegisterTMDBRouter(tmdbRouter *gin.RouterGroup) {
	imgRouter := tmdbRouter.Group("/image") // 图片相关接口
	{
		imgRouter.GET("/poster/:media_type/:tmdb_id", PosterImage) // 获取媒体海报图片
	}
	tmdbRouter.GET("/overview/:media_type/:tmdb_id", Overview) // 获取概述

}
