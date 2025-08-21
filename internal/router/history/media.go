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
)

// @Router /history/media/transfer [get]
// @Summary 查询媒体转移历史记录
// @Description 查询媒体转移历史记录，支持根据 ID、时间范围、源路径、目标路径、转移类型和状态进行过滤
// @Tags History
// @Produce json
// @Param id query uint false "媒体转移历史记录 ID, 如果提供则只查询该 ID 的记录"
// @Param start_time query time.Time false "开始时间, 格式为 RFC3339"
// @Param end_time query time.Time false "结束时间, 格式为 RFC3339"
// @Param storage_type query string false "存储类型, 可选值为 'LocalStorage' 等"
// @Param path query string false "路径, 模糊匹配"
// @Param transfer_type query string false "转移类型, 可选值为 'Copy'、'Move'、'Link'、'SoftLink' 等"
// @Param status query bool false "是否成功, true 或 false"
// @Param count query int false "最大返回数量, 默认值为 30"
func QueryMediaTransferHistory(ctx *gin.Context) {
	var (
		resp schemas.Response[[]*models.MediaTransferHistory]

		id                 *uint                // ID
		startTime, endTime *time.Time           // 时间范围
		storageType        storage.StorageType  // 存储类型
		path               string               // 路径，模糊匹配
		transferType       storage.TransferType // 转移类型
		status             *bool                // 是否成功
		count              = 30                 // 默认返回数量
	)

	// 解析 ID
	idStr := ctx.Query("id")
	if idStr != "" {
		idp, err := strconv.Atoi(idStr)
		if err != nil {
			resp.Message = "解析 ID 参数失败: " + err.Error()
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		idUint := uint(idp) // 先转换为 uint
		id = &idUint        // 然后取地址赋给指针
	} else { // 如果没有提供 ID，则根据其他条件查询
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
	}

	var respHistories []*models.MediaTransferHistory
	for history, err := range database.QueryMediaTransferHistory(ctx, id, startTime, endTime, storageType, path, transferType, status) {
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
