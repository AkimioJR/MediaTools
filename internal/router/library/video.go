package library

import (
	"MediaTools/internal/controller/library_controller"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ArchiveMediaManual(ctx *gin.Context) {
	var (
		req  schemas.ArchiveMediaManualRequest
		resp schemas.Response[*storage.StorageFileInfo]
	)

	req.Season = -1 // 默认值为 -1，表示不设定季编号
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}
	srcFile := storage.NewStoragePath(req.SrcFile.StorageType, req.SrcFile.Path)
	dstDir := storage.NewStoragePath(req.DstDir.StorageType, req.DstDir.Path)

	dstFile, err := library_controller.ArchiveMediaAdvanced(srcFile, dstDir, req.TransferType, req.MediaType,
		req.TMDBID, req.Season, req.EpisodeStr, req.EpisodeOffset, req.Part,
		req.OrganizeByType, req.OrganizeByCategory, req.Scrape,
	)
	if err != nil {
		resp.Message = "整理媒体文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	logrus.Info("手动整理媒体文件整理成功：", dstFile.String())
	resp.RespondSuccessJSON(ctx, dstFile.(*storage.StorageFileInfo))
}
