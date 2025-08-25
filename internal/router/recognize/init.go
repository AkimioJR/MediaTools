package recognize

import "github.com/gin-gonic/gin"

func RegisteRecognizeRouter(recognizeRouter *gin.RouterGroup) {
	recognizeRouter.GET("/media", RecognizeMedia) // 识别媒体信息
}
