package models

import (
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
)

// 视频媒体转移历史记录
type MediaTransferHistory struct {
	BaseModel
	SrcType      storage.StorageType  `json:"src_type"`              // 源存存储器
	SrcPath      string               `json:"src_path"`              // 源路径
	DstType      storage.StorageType  `json:"dst_type"`              // 目标存储器
	DstPath      string               `json:"dst_path"`              // 目标路径
	TransferType storage.TransferType `json:"transfer_type"`         // 转移类型
	Status       bool                 `json:"status"`                // 是否成功
	Message      string               `json:"message"`               // 错误信息
	Item         *schemas.MediaItem   `json:"item" gorm:"type:text"` // 识别到的媒体信息
}
