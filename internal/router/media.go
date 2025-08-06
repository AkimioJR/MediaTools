package router

import (
	"MediaTools/internal/controller/media_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /media
// @Route /recognize [get]
// @Summary 识别媒体信息
// @Description 根据提供的标题识别媒体信息，并返回 MediaItem 对象
// @Tags media
// @Param title query string true "媒体标题"
// @Success 200 {object} Response[*schemas.MediaItem]
// @Failure 400 {object} Response[*schemas.MediaItem]
// @Failure 500 {object} Response[*schemas.MediaItem]
func MediaRecognize(ctx *gin.Context) {
	var resp Response[*schemas.MediaItem]
	title := ctx.Query("title")
	if title == "" {
		resp.Message = "标题不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	logrus.Infof("正在识别媒体：%s", title)
	videoMeta := media_controller.ParseVideoMeta(title)
	mediaInfo, err := tmdb_controller.RecognizeAndEnrichMedia(videoMeta, nil, nil)
	if err != nil {
		resp.Message = "识别失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	item, err := schemas.NewMediaItem(videoMeta, mediaInfo)
	if err != nil {
		resp.Message = "创建媒体项失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Data = item
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}
