package router

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
)

type PathRequest struct {
	Path string `json:"path" binding:"required"`
}
type FileInfoRequest struct {
	StorageType schemas.StorageType `json:"storage_type" binding:"required"`
	Path        string              `json:"path" binding:"required"`
}

type TransferRequest struct {
	SrcFile      FileInfoRequest      `json:"src_file" binding:"required"`
	DstFile      FileInfoRequest      `json:"dst_file" binding:"required"`
	TransferType schemas.TransferType `json:"transfer_type"`
}

type ScrapeRequest struct {
	DstFile   FileInfoRequest `json:"dst_file" binding:"required"`
	MediaType *meta.MediaType `json:"media_type,omitempty"`
	TMDBID    *int            `json:"tmdb_id,omitempty"`
}

type LibraryArchiveMediaRequest struct {
	SrcFile      FileInfoRequest      `json:"src_file" binding:"required"`
	DstDir       FileInfoRequest      `json:"dst_dir" binding:"required"`
	TransferType schemas.TransferType `json:"transfer_type"`
	NeedScrape   bool                 `json:"need_scrape"`
}
