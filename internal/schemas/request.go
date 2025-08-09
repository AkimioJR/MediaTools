package schemas

import (
	"MediaTools/internal/pkg/meta"
)

type PathRequest struct {
	Path string `json:"path" binding:"required"`
}
type FileInfoRequest struct {
	StorageType StorageType `json:"storage_type" binding:"required"`
	Path        string      `json:"path" binding:"required"`
}

type TransferRequest struct {
	SrcFile      FileInfoRequest `json:"src_file" binding:"required"`
	DstFile      FileInfoRequest `json:"dst_file" binding:"required"`
	TransferType TransferType    `json:"transfer_type"`
}

type ScrapeRequest struct {
	DstFile   FileInfoRequest `json:"dst_file" binding:"required"`
	MediaType *meta.MediaType `json:"media_type,omitempty"`
	TMDBID    *int            `json:"tmdb_id,omitempty"`
}

type LibraryArchiveMediaRequest struct {
	SrcFile      FileInfoRequest `json:"src_file" binding:"required"`
	DstDir       FileInfoRequest `json:"dst_dir" binding:"required"`
	TransferType TransferType    `json:"transfer_type"`
	NeedScrape   bool            `json:"need_scrape"`
}

type RenameRequest struct {
	Path    string `json:"path" binding:"required"`
	NewName string `json:"new_name" binding:"required"`
}
