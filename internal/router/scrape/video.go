package scrape

import (
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /scrape
// @Router /video [post]
// @Summary 刮削视频
// @Description 刮削视频文件的元数据和相关信息
// @Tags scrape
// @Accept json
// @Produce json
// @Param request body schemas.ScrapeRequest true "刮削请求参数"
// @Success 200 {object} schemas.Response[*schemas.FileInfo] "刮削成功"
// @Failure 400 {object} schemas.Response[*schemas.FileInfo] "请求参数错误"
// @Failure 500 {object} schemas.Response[*schemas.FileInfo] "刮削失败
func Video(ctx *gin.Context) {
	var (
		req  schemas.ScrapeRequest
		resp schemas.Response[*schemas.FileInfo]
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			resp.Message = "请求参数错误: " + err.Error()
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
	}

	dstFile := schemas.FileInfo{
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

	err := scrape_controller.RecognizeAndScrape(&dstFile, req.MediaType, req.TMDBID)
	if err != nil {
		resp.Message = "刮削失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Success = true
	resp.Data = &dstFile
	ctx.JSON(http.StatusOK, resp)
}
