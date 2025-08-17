package schemas

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas/storage"
)

type PathRequest struct {
	Path string `json:"path" binding:"required"`
}
type FileInfoRequest struct {
	StorageType storage.StorageType `json:"storage_type" binding:"required"`
	Path        string              `json:"path" binding:"required"`
}

type TransferRequest struct {
	SrcFile      FileInfoRequest      `json:"src_file" binding:"required"`
	DstFile      FileInfoRequest      `json:"dst_file" binding:"required"`
	TransferType storage.TransferType `json:"transfer_type"`
}

type ScrapeRequest struct {
	DstFile   FileInfoRequest `json:"dst_file" binding:"required"`
	MediaType meta.MediaType  `json:"media_type,omitempty"`
	TMDBID    int             `json:"tmdb_id,omitempty"`
}

type LibraryArchiveMediaRequest struct {
	SrcFile      FileInfoRequest      `json:"src_file" binding:"required"`
	DstDir       FileInfoRequest      `json:"dst_dir" binding:"required"`
	TransferType storage.TransferType `json:"transfer_type"`
	Scrape       bool                 `json:"scrape"`
}

type RenameRequest struct {
	Path    string `json:"path" binding:"required"`
	NewName string `json:"new_name" binding:"required"`
}

type ArchiveMediaManualRequest struct {
	SrcFile            FileInfoRequest      `json:"src_file" binding:"required"`             // 源文件
	DstDir             FileInfoRequest      `json:"dst_dir" binding:"required"`              // 目标目录
	TransferType       storage.TransferType `json:"transfer_type" binding:"required"`        // 转移方法
	OrganizeByType     bool                 `json:"organize_by_type" binding:"required"`     // 是否按类型整理
	OrganizeByCategory bool                 `json:"organize_by_category" binding:"required"` // 是否按分类整理
	Scrape             bool                 `json:"scrape" binding:"required"`               // 是否刮削元数据

	// 可选字段
	MediaType     meta.MediaType `json:"media_type"`     // 媒体类型
	TMDBID        int            `json:"tmdb_id"`        // TMDB ID
	Season        int            `json:"season"`         // 季编号，-1 表示不设定
	EpisodeStr    string         `json:"episode_str"`    // 集数字符串，单集或多集范围
	EpisodeOffset string         `json:"episode_offset"` // 集数偏移（仅为制定集数是生效）
	Part          string         `json:"part"`           // 指定分段
}
