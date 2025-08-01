package router

import "MediaTools/internal/schemas"

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
