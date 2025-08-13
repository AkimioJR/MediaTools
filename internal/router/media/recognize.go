package media

import (
	"MediaTools/internal/controller/media_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Route /media/recognize [get]
// @Summary 识别媒体信息
// @Description 根据提供的标题识别媒体信息，并返回 MediaItem 对象
// @Tags 媒体信息
// @Param title query string true "媒体标题"
// @Produce json
// @Success 200 {object} schemas.MediaItem
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func Recognize(ctx *gin.Context) {
	var errResp schemas.ErrResponse
	title := ctx.Query("title")
	if title == "" {
		errResp.Message = "标题不能为空"
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}
	logrus.Infof("正在识别媒体：%s", title)
	videoMeta, _ := media_controller.ParseVideoMeta(title)
	mediaInfo, err := tmdb_controller.RecognizeAndEnrichMedia(videoMeta)
	if err != nil {
		errResp.Message = "识别失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}
	item, err := schemas.NewMediaItem(videoMeta, mediaInfo)
	if err != nil {
		errResp.Message = "创建媒体项失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}
	ctx.JSON(http.StatusOK, item)
}
