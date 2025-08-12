package library

import (
	"MediaTools/internal/controller/library_controller"
	"MediaTools/internal/controller/media_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"
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
// @Param request body LibraryArchiveMediaRequest true "请求参数"
// @Success 200 {object} schemas.Response[schemas.FileInfo] "成功响应"
// @Failure 400 {object} schemas.Response[schemas.FileInfo] "请求参数错误"
// @Failure 500 {object} schemas.Response[schemas.FileInfo] "服务器错误"
func LibraryArchiveMedia(ctx *gin.Context) {
	var (
		req  schemas.LibraryArchiveMediaRequest
		resp schemas.Response[*schemas.FileInfo]
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	srcFile, err := storage_controller.GetFile(req.SrcFile.Path, req.SrcFile.StorageType)
	if err != nil {
		resp.Message = fmt.Sprintf("获取源文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	dstDir, err := storage_controller.GetFile(req.DstDir.Path, req.DstDir.StorageType)
	if err != nil {
		resp.Message = fmt.Sprintf("获取目标目录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	logrus.Info("正在解析视频元数据：", srcFile.Name)
	videoMeta := media_controller.ParseVideoMeta(srcFile.Name)
	info, err := tmdb_controller.RecognizeMedia(videoMeta)
	if err != nil {
		resp.Message = "识别媒体信息失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	if info.MediaType == meta.MediaTypeTV {
		seasonDetail, err := tmdb_controller.GetTVSeasonDetail(info.TMDBID, videoMeta.Season)
		if err != nil {
			resp.Message = "获取电视剧季节信息失败: " + err.Error()
			ctx.JSON(http.StatusInternalServerError, resp)
			return
		}
		info.TMDBInfo.TVInfo.SeasonInfo = seasonDetail.TMDBInfo.TVInfo.SeasonInfo
		episodeDetail, err := tmdb_controller.GetTVEpisodeDetail(info.TMDBID, videoMeta.Season, videoMeta.Episode)
		if err != nil {
			resp.Message = "获取电视剧集信息失败: " + err.Error()
			ctx.JSON(http.StatusInternalServerError, resp)
			return
		}
		info.TMDBInfo.TVInfo.EpisodeInfo = episodeDetail.TMDBInfo.TVInfo.EpisodeInfo
	}
	item, err := schemas.NewMediaItem(videoMeta, info)
	if err != nil {
		resp.Message = "创建媒体项失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	var dst *schemas.FileInfo
	if req.NeedScrape {
		dst, err = library_controller.ArchiveMedia(srcFile, dstDir, schemas.TransferLink, item, info)
	} else {
		dst, err = library_controller.ArchiveMedia(srcFile, dstDir, schemas.TransferLink, item, nil)
	}
	if err != nil {
		resp.Message = "转移媒体文件失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	logrus.Infof("%s 媒体文件转移完成", srcFile)

	resp.Data = dst
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}
