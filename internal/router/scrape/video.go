package scrape

import (
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Router /scrape/video [post]
// @Summary 刮削视频
// @Description 刮削视频文件的元数据和相关信息
// @Tags 刮削
// @Accept json
// @Produce json
// @Param request body schemas.ScrapeRequest true "刮削请求参数"
func Video(ctx *gin.Context) {
	var (
		req  schemas.ScrapeRequest
		resp schemas.Response[*storage.StorageFileInfo]
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	dstFile := storage.StorageFileInfo{
		StorageType: req.DstFile.StorageType,
		Path:        req.DstFile.Path,
	}

	switch {
	case req.MediaType != meta.MediaTypeUnknown && req.TMDBID != 0:
		logrus.Infof("开始刮削视频：%s，媒体类型：%s，TMDB ID：%d", dstFile.String(), req.MediaType, req.TMDBID)
	case req.MediaType != meta.MediaTypeUnknown:
		logrus.Infof("开始刮削视频：%s，媒体类型：%s", dstFile.String(), req.MediaType)
	case req.TMDBID != 0:
		logrus.Infof("开始刮削视频：%s，TMDB ID：%d", dstFile.String(), req.TMDBID)
	default:
		logrus.Infof("开始刮削视频：%s", dstFile.String())
	}

	err := scrape_controller.RecognizeAndScrape(ctx, &dstFile, req.MediaType, req.TMDBID)
	if err != nil {
		resp.Message = "刮削失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Data = &dstFile
	resp.RespondJSON(ctx, http.StatusOK)
}
