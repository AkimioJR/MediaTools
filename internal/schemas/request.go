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
