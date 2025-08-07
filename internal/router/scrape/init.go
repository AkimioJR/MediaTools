package scrape

import "github.com/gin-gonic/gin"

func RegisterScrapeRouter(router *gin.Engine) {
	scrapeRouter := router.Group("/scrape") // 刮削相关接口
	{
		scrapeRouter.POST("/video", Video) // 刮削视频
	}
}
