package recognize

import (
	"MediaTools/internal/controller/recognize_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Route /recognize/media [get]
// @Summary 识别媒体信息
// @Description 根据提供的标题识别媒体信息，并返回 MediaItem 对象
// @Tags 识别
// @Param title query string true "媒体标题"
// @Produce json
func RecognizeMedia(ctx *gin.Context) {
	var resp schemas.Response[*schemas.RecognizeMediaDetail]

	title := ctx.Query("title")
	if title == "" {
		resp.Message = "标题不能为空"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	logrus.Infof("正在识别媒体：%s", title)
	videoMeta, customRule, metaRule := recognize_controller.ParseVideoMeta(title)
	mediaInfo, err := tmdb_controller.RecognizeAndEnrichMedia(ctx,videoMeta)
	if err != nil {
		resp.Message = "识别失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	item, err := schemas.NewMediaItem(videoMeta, mediaInfo)
	if err != nil {
		resp.Message = "创建媒体项失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Data = &schemas.RecognizeMediaDetail{
		Item:       item,
		CustomRule: customRule,
		MetaRule:   metaRule,
	}
	resp.RespondJSON(ctx, http.StatusOK)
}
