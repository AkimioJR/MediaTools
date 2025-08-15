package library

import (
	"MediaTools/internal/controller/library_controller"
	"MediaTools/internal/controller/recognize_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Router  /library/archive [post]
// @Summary 归档媒体文件
// @Description 整理一个视频文件及其相关的字幕和音轨文件到指定目录
// @Tags 媒体库
// @Accept json
// @Produce json
// @Param request body schemas.LibraryArchiveMediaRequest true "请求参数"
func LibraryArchiveMedia(ctx *gin.Context) {
	var (
		req  schemas.LibraryArchiveMediaRequest
		resp schemas.Response[*storage.StorageFileInfo]
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	srcPath, err := storage_controller.GetPath(req.SrcFile.Path, req.SrcFile.StorageType)
	if err != nil {
		resp.Message = "获取源文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	dstDir, err := storage_controller.GetPath(req.DstDir.Path, req.DstDir.StorageType)
	if err != nil {
		resp.Message = "获取目标目录失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	logrus.Infof("正在解析视频元数据：%s", srcPath.String())
	videoMeta, _, _ := recognize_controller.ParseVideoMeta(filepath.Base(srcPath.GetPath()))
	info, err := tmdb_controller.RecognizeMedia(videoMeta)
	if err != nil {
		resp.Message = "识别媒体信息失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	if info.MediaType == meta.MediaTypeTV {
		seasonDetail, err := tmdb_controller.GetTVSeasonDetail(info.TMDBID, videoMeta.Season)
		if err != nil {
			resp.Message = "获取电视剧季节信息失败: " + err.Error()
			resp.RespondJSON(ctx, http.StatusInternalServerError)
			return
		}
		info.TMDBInfo.TVInfo.SeasonInfo = seasonDetail.TMDBInfo.TVInfo.SeasonInfo
		episodeDetail, err := tmdb_controller.GetTVEpisodeDetail(info.TMDBID, videoMeta.Season, videoMeta.Episode)
		if err != nil {
			resp.Message = "获取电视剧集信息失败: " + err.Error()
			resp.RespondJSON(ctx, http.StatusInternalServerError)
			return
		}
		info.TMDBInfo.TVInfo.EpisodeInfo = episodeDetail.TMDBInfo.TVInfo.EpisodeInfo
	}

	item, err := schemas.NewMediaItem(videoMeta, info)
	if err != nil {
		resp.Message = "创建媒体项失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	var dst storage.StoragePath
	if req.NeedScrape {
		dst, err = library_controller.ArchiveMedia(srcPath, dstDir, storage.TransferLink, item, info)
	} else {
		dst, err = library_controller.ArchiveMedia(srcPath, dstDir, storage.TransferLink, item, nil)
	}
	if err != nil {
		resp.Message = "转移媒体文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}
	logrus.Infof("%s 媒体文件转移完成", srcPath)

	fileInfo, err := storage_controller.GetDetail(dst)
	if err != nil {
		resp.Message = "获取目标文件信息失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}
	resp.Data = fileInfo
	resp.RespondJSON(ctx, http.StatusOK)
}
