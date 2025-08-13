package library

import (
	"MediaTools/internal/controller/library_controller"
	"MediaTools/internal/controller/media_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /library
// @Router /archive [post]
// @Summary 归档媒体文件
// @Description 整理一个视频文件及其相关的字幕和音轨文件到指定目录
// @Tags Library
// @Accept json
// @Produce json
// @Param request body schemas.LibraryArchiveMediaRequest true "请求参数"
// @Success 200 {object} schemas.FileInfo "成功响应"
// @Failure 400 {object} schemas.ErrResponse "请求参数错误"
// @Failure 500 {object} schemas.ErrResponse "服务器错误"
func LibraryArchiveMedia(ctx *gin.Context) {
	var (
		req     schemas.LibraryArchiveMediaRequest
		errResp schemas.ErrResponse
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errResp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	srcFile, err := storage_controller.GetFile(req.SrcFile.Path, req.SrcFile.StorageType)
	if err != nil {
		errResp.Message = "获取源文件失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}
	dstDir, err := storage_controller.GetFile(req.DstDir.Path, req.DstDir.StorageType)
	if err != nil {
		errResp.Message = "获取目标目录失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	logrus.Info("正在解析视频元数据：", srcFile.Name)
	videoMeta := media_controller.ParseVideoMeta(srcFile.Name)
	info, err := tmdb_controller.RecognizeMedia(videoMeta)
	if err != nil {
		errResp.Message = "识别媒体信息失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}
	if info.MediaType == meta.MediaTypeTV {
		seasonDetail, err := tmdb_controller.GetTVSeasonDetail(info.TMDBID, videoMeta.Season)
		if err != nil {
			errResp.Message = "获取电视剧季节信息失败: " + err.Error()
			ctx.JSON(http.StatusInternalServerError, errResp)
			return
		}
		info.TMDBInfo.TVInfo.SeasonInfo = seasonDetail.TMDBInfo.TVInfo.SeasonInfo
		episodeDetail, err := tmdb_controller.GetTVEpisodeDetail(info.TMDBID, videoMeta.Season, videoMeta.Episode)
		if err != nil {
			errResp.Message = "获取电视剧集信息失败: " + err.Error()
			ctx.JSON(http.StatusInternalServerError, errResp)
			return
		}
		info.TMDBInfo.TVInfo.EpisodeInfo = episodeDetail.TMDBInfo.TVInfo.EpisodeInfo
	}
	item, err := schemas.NewMediaItem(videoMeta, info)
	if err != nil {
		errResp.Message = "创建媒体项失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}
	var dst *schemas.FileInfo
	if req.NeedScrape {
		dst, err = library_controller.ArchiveMedia(srcFile, dstDir, schemas.TransferLink, item, info)
	} else {
		dst, err = library_controller.ArchiveMedia(srcFile, dstDir, schemas.TransferLink, item, nil)
	}
	if err != nil {
		errResp.Message = "转移媒体文件失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}
	logrus.Infof("%s 媒体文件转移完成", srcFile)
	ctx.JSON(http.StatusOK, dst)
}
