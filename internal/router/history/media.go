package history

import (
	"MediaTools/internal/database"
	"MediaTools/internal/models"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Router /history/media/transfer [get]
// @Summary 查询媒体转移历史记录
// @Description 查询媒体转移历史记录，支持根据 ID、时间范围、源路径、目标路径、转移类型和状态进行过滤
// @Tags 历史记录
// @Produce json
// @Param start_time query time.Time false "开始时间, 格式为 RFC3339"
// @Param end_time query time.Time false "结束时间, 格式为 RFC3339"
// @Param storage_type query string false "存储类型, 可选值为 'LocalStorage' 等"
// @Param path query string false "路径, 模糊匹配"
// @Param transfer_type query string false "转移类型, 可选值为 'Copy'、'Move'、'Link'、'SoftLink' 等"
// @Param status query bool false "是否成功, true 或 false"
// @Param count query int false "最大返回数量, 默认值为 50"
// @Param page query int false "页码, 从 1 开始, 默认值为 1"
func QueryMediaTransferHistory(ctx *gin.Context) {
	var (
		resp schemas.Response[[]*models.MediaTransferHistory]

		startTime, endTime *time.Time           // 时间范围
		storageType        storage.StorageType  // 存储类型
		path               string               // 路径，模糊匹配
		transferType       storage.TransferType // 转移类型
		status             *bool                // 是否成功
		count              = 50                 // 默认返回数量
		page               = 1                  // 默认页码
	)

	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	if startTimeStr != "" {
		t, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			resp.Message = "解析开始时间失败: " + err.Error()
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		startTime = &t
	}
	if endTimeStr != "" {
		t, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			resp.Message = "解析结束时间失败: " + err.Error()
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		endTime = &t
	}
	storageTypeStr := ctx.Query("storage_type")
	storageType = storage.ParseStorageType(storageTypeStr)
	path = ctx.Query("path")

	transferTypeStr := ctx.Query("transfer_type")
	transferType = storage.ParseTransferType(transferTypeStr)

	statusStr := ctx.Query("status")

	if statusStr != "" {
		b, err := strconv.ParseBool(statusStr)
		if err != nil {
			resp.Message = "解析状态参数失败: " + err.Error()
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		status = &b
	}
	countStr := ctx.Query("count")
	if countStr != "" {
		c, err := strconv.Atoi(countStr)
		if err != nil {
			resp.Message = "解析数量参数失败: " + err.Error()
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		count = c
	}
	pageStr := ctx.Query("page")
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			resp.Message = "解析页码参数失败: " + err.Error()
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		if p < 1 {
			logrus.Warningf("页码参数无效，重置为 1: %d", p)
			p = 1
		}
		page = p
	}

	offset := (page - 1) * count
	respHistories := make([]*models.MediaTransferHistory, 0, count)
	histories, err := database.QueryMediaTransferHistory(ctx, startTime, endTime, storageType, path, transferType, status, offset)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	for history, err := range histories {
		if err != nil {
			resp.Message = err.Error()
			resp.RespondJSON(ctx, http.StatusInternalServerError)
			return
		}

		if count == 0 {
			break
		}
		count--
		respHistories = append(respHistories, history)
	}
	resp.RespondSuccessJSON(ctx, respHistories)
}

// @Router /media/{id} [get]
// @Summary 查询媒体转移历史记录 by ID
// @Description 根据 ID 查询媒体转移历史记录
// @Tag 历史记录
// @Produce json
// @Param id path uint64 true "媒体转移历史记录 ID"
func QueryMediaTransferHistoryByID(ctx *gin.Context) {
	var resp schemas.Response[*models.MediaTransferHistory]
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp.Message = "无效的 ID 参数: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	history, err := database.QueryMediaTransferHistoryByID(ctx, id)
	if err != nil {
		resp.Message = "查询媒体转移历史记录失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.RespondSuccessJSON(ctx, history)
}

// @Router /history/media/{id} [delete]
// @Summary 删除媒体转移历史记录
// @Description 根据 ID 删除媒体转移历史记录
// @Tag 历史记录
// @Param id path uint64 true "媒体转移历史记录 ID"
func DeleteMediaTransferHistory(ctx *gin.Context) {
	var resp schemas.Response[any]

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp.Message = "无效的 ID 参数: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	err = database.DeleteMediaTransferHistory(ctx, id)
	if err != nil {
		resp.Message = "删除媒体转移历史记录失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.RespondSuccessJSON(ctx, nil)
}
