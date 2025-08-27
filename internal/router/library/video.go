package library

import (
	"MediaTools/internal/controller/library_controller"
	"MediaTools/internal/pkg/task"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Router /library/archive [post]
// @Summary 手动归档媒体文件
// @Description 手动归档媒体文件
// @Tags 媒体库管理
// @Accept json
// @Produce json
// @Param data body schemas.ArchiveMediaManualRequest true "请求参数"
func ArchiveMediaManual(ctx *gin.Context) {
	var (
		req  schemas.ArchiveMediaManualRequest
		resp schemas.Response[*task.Task]
	)

	req.Season = -1 // 默认值为 -1，表示不设定季编号
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}
	srcFile := storage.NewStoragePath(req.SrcFile.StorageType, req.SrcFile.Path)
	dstDir := storage.NewStoragePath(req.DstDir.StorageType, req.DstDir.Path)

	task, err := library_controller.ArchiveMediaAdvanced(ctx, srcFile, dstDir, req.TransferType, req.MediaType,
		req.TMDBID, req.Season, req.EpisodeStr, req.EpisodeFormat, req.EpisodeOffset, req.Part,
		req.OrganizeByType, req.OrganizeByCategory, req.Scrape,
	)
	if err != nil {
		resp.Message = "整理媒体文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	logrus.Debugf("已提交媒体文件整理任务: %+v", task)
	resp.RespondSuccessJSON(ctx, task)
}
