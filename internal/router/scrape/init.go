package scrape

import "github.com/gin-gonic/gin"

func RegisterScrapeRouter(scrapeRouter *gin.RouterGroup) {
	scrapeRouter.POST("/video", Video) // 刮削视频
}
