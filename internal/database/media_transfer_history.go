package database

import (
	"MediaTools/internal/models"
	"MediaTools/internal/schemas/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UpdateMediaTransferHistory 更新媒体转移历史记录
// 如果 ID 不存在则创建新的转移历史记录
// 如果 ID 存在则更新现有的转移历史记录
func UpdateMediaTransferHistory(history *models.MediaTransferHistory) error {
	result := DB.Save(history)
	if result.Error != nil {
		return fmt.Errorf("更新媒体转移历史记录失败: %w", result.Error)
	}
	return nil
}

func QueryMediaTransferHistoryBySrc(src storage.StoragePath) (*models.MediaTransferHistory, error) {
	ctx := context.Background()
	history, err := gorm.G[models.MediaTransferHistory](DB).Where("src_type = ? AND src_path = ?", src.GetStorageType(), src.GetPath()).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询媒体转移历史记录失败: %w", err)
	}
	return &history, nil
}

// QueryMediaTransferHistory 查询媒体转移历史记录
// 支持根据 ID、时间范围、源路径、目标路径、转移类型和状态进行过滤
// 如果 ID 不为 nil，则只查询该 ID 的记录
// 如果 ID 为 nil，则根据其他条件查询
func QueryMediaTransferHistory(
	ctx context.Context,
	id *uint, startTime *time.Time, endTime *time.Time,
	storageType storage.StorageType, // 存储类型和路径，存储类型为 StorageUnknown 时不进行过滤，否则对 src 和 dst 都进行过滤
	path string, // 路径，模糊匹配
	transferType storage.TransferType, // 转移类型为 TransferUnknown 时不进行过滤
	status *bool, // 是否成功
	count int, // 最大返回数量
) ([]models.MediaTransferHistory, error) {
	var query gorm.ChainInterface[models.MediaTransferHistory] = gorm.G[models.MediaTransferHistory](DB)

	if id != nil { // 如果提供了 ID，则只查询该 ID 的记录
		query = query.Where("id = ?", *id)
	} else { // 如果没有提供 ID，则根据其他条件查询

		switch {
		case startTime != nil && endTime != nil: // 如果同时提供了开始和结束时间
			query = query.Where("created_at BETWEEN ? AND ?", *startTime, *endTime)
		case startTime != nil: // 如果只提供了开始时间
			query = query.Where("created_at >= ?", *startTime)
		case endTime != nil: // 如果只提供了结束时间
			query = query.Where("created_at <= ?", *endTime)
		}

		if storageType != storage.StorageUnknown {
			query = query.Where("src_type = ? OR dst_type = ?", storageType, storageType)
		}
		if path != "" {
			query = query.Where("src_path LIKE ? OR dst_path LIKE ?", "%"+path+"%", "%"+path+"%")
		}

		if transferType != storage.TransferUnknown {
			query = query.Where("transfer_type = ?", transferType)
		}

		if status != nil {
			query = query.Where("status = ?", *status)
		}
		if count > 0 {
			query = query.Limit(count)
		}
	}

	results, err := query.Find(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return results, nil // 如果没有找到记录，返回空切片
		} else {
			return nil, fmt.Errorf("查询媒体转移历史记录失败: %w", err)
		}
	}
	return results, nil
}
