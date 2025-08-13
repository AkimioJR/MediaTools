package recognize

import "github.com/gin-gonic/gin"

func RegisteRecognizeRouter(router *gin.Engine) {
	recognizeRouter := router.Group("/recognize") // 媒体相关接口
	{
		recognizeRouter.GET("/media", RecognizeMedia) // 识别媒体信息
	}
}
