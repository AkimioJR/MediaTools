package tmdb

import (
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @BasePath /tmdb
// @Router /overview/{media_type}/{tmdb_id} [get]
// @Summary 获取媒体概述信息
// @Description 根据媒体类型和 TMDB ID 获取对应的概述信息
// @Tags TMDB
// @Param media_type path string true "媒体类型" Enums(Movie, TV)
// @Param tmdb_id path uint true "TMDB ID"
// @Param season query uint false "季数"
// @Param episode query uint false "集数"
func Overview(ctx *gin.Context) {
	var errResp schemas.ErrResponse

	// 获取类型和 ID
	mediaTypeStr := ctx.Param("media_type")
	tmdbIDStr := ctx.Param("tmdb_id")
	tmdbID, err := strconv.Atoi(tmdbIDStr)
	if err != nil {
		errResp.Message = "非法 TMDB ID: " + err.Error()
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}
	var (
		season  = -1
		episode = -1
	)
	seasonStr := ctx.Query("season")
	episodeStr := ctx.Query("episode")
	if seasonStr != "" {
		season, err = strconv.Atoi(seasonStr)
		if err != nil {
			errResp.Message = "非法季数参数: " + err.Error()
			ctx.JSON(http.StatusBadRequest, errResp)
			return
		}
	}
	if episodeStr != "" {
		episode, err = strconv.Atoi(episodeStr)
		if err != nil {
			errResp.Message = "非法集数参数: " + err.Error()
			ctx.JSON(http.StatusBadRequest, errResp)
			return
		}
	}

	var overview string
	switch meta.ParseMediaType(mediaTypeStr) {
	case meta.MediaTypeUnknown:
		errResp.Message = "无效的媒体类型"
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	case meta.MediaTypeMovie: // 处理电影类型
		movieInfo, err := tmdb_controller.GetMovieDetail(tmdbID)
		if err != nil {
			errResp.Message = "获取电影详情失败: " + err.Error()
			ctx.JSON(http.StatusInternalServerError, errResp)
			return
		}
		overview = movieInfo.TMDBInfo.MovieInfo.Overview
	case meta.MediaTypeTV: // 处理电视剧类型
		switch {
		case season >= 0 && episode > 0: // 获取特定季集的概述
			episodeInfo, err := tmdb_controller.GetTVEpisodeDetail(tmdbID, season, episode)
			if err != nil {
				errResp.Message = "获取电视剧集详情失败: " + err.Error()
				ctx.JSON(http.StatusInternalServerError, errResp)
				return
			}
			overview = episodeInfo.TMDBInfo.TVInfo.EpisodeInfo.Overview
		case season >= 0: // 获取特定季的概述
			seasonInfo, err := tmdb_controller.GetTVSeasonDetail(tmdbID, season)
			if err != nil {
				errResp.Message = "获取电视剧季详情失败: " + err.Error()
				ctx.JSON(http.StatusInternalServerError, errResp)
				return
			}
			overview = seasonInfo.TMDBInfo.TVInfo.SeasonInfo.Overview
		default: // 获取整部剧的概述
			tvInfo, err := tmdb_controller.GetTVSerieDetail(tmdbID)
			if err != nil {
				errResp.Message = "获取电视剧详情失败: " + err.Error()
				ctx.JSON(http.StatusInternalServerError, errResp)
				return
			}
			overview = tvInfo.TMDBInfo.TVInfo.SerieInfo.Overview
		}
	}

	if overview == "" {
		errResp.Message = "未找到媒体概述信息"
		ctx.JSON(http.StatusNotFound, errResp)
		return
	}
	ctx.JSON(http.StatusOK, overview)
}
