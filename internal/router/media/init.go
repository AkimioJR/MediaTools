package media

import "github.com/gin-gonic/gin"

func RegisterMediaRouter(router *gin.Engine) {
	mediaRouter := router.Group("/media") // 媒体相关接口
	{
		mediaRouter.GET("/recognize", Recognize)
	}
}
