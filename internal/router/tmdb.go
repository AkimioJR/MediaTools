package router

import (
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @BasePath /tmdb/image
// @Router /poster/{media_type}/{tmdb_id} [get]
// @Summary 获取媒体海报图片
// @Description 根据媒体类型和 TMDB ID 获取对应的海报图片 URL
// @Tags TMDB
// @Param media_type path string true "媒体类型" Enums(Movie, TV)
// @Param tmdb_id path int true "TMDB ID"
// @Success 200 {object} Response[*string] "成功响应，返回图片 URL"
// @Failure 400 {object} Response[*string] "请求参数错误"
// @Failure 500 {object} Response[*string] "服务器错误"
func TMDBPosterImage(ctx *gin.Context) {
	var resp Response[*string]

	// 获取类型和 ID
	mediaTypeStr := ctx.Param("media_type")
	tmdbIDStr := ctx.Param("tmdb_id")
	tmdbID, err := strconv.Atoi(tmdbIDStr)
	if err != nil {
		resp.Message = "非法 TMDB ID: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var imgPath string
	switch meta.ParseMediaType(mediaTypeStr) {
	case meta.MediaTypeUnknown:
		resp.Message = "无效的媒体类型"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	case meta.MediaTypeMovie: // 处理电影类型
		imagesInfo, err := tmdb_controller.GetMovieImage(tmdbID)
		if err != nil {
			resp.Message = "获取电影图片信息失败: " + err.Error()
			ctx.JSON(http.StatusInternalServerError, resp)
			return
		}
		if len(imagesInfo.Posters) == 0 {
			resp.Message = "未找到电影海报图片"
			ctx.JSON(http.StatusNotFound, resp)
			return
		}
		imgPath = imagesInfo.Posters[0].FilePath
	case meta.MediaTypeTV: // 处理电视剧类型
		imagesInfo, err := tmdb_controller.GetTVSerieImage(tmdbID)
		if err != nil {
			resp.Message = "获取电视剧详情失败: " + err.Error()
			ctx.JSON(http.StatusInternalServerError, resp)
			return
		}
		if len(imagesInfo.Posters) == 0 {
			resp.Message = "未找到电视剧海报图片"
			ctx.JSON(http.StatusNotFound, resp)
			return
		}
		imgPath = imagesInfo.Posters[0].FilePath
	}
	imgURL := tmdb_controller.GetImageURL(imgPath)
	resp.Data = &imgURL
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}
