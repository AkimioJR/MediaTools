package storage

import (
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	StorageType StorageType `json:"storage_type"` // 存储系统类型
	Path        string      `json:"path"`         // 文件路径
	Name        string      `json:"name"`         // 文件名
	Ext         string      `json:"ext"`          // 文件扩展名
	Size        int64       `json:"size"`         // 文件大小
	IsDir       bool        `json:"is_dir"`       // 是否为目录
	ModTime     time.Time   `json:"mod_time"`     // 文件修改时间
}

func NewBasicFileInfo(storageType StorageType, path string) *FileInfo {
	return &FileInfo{
		StorageType: storageType,
		Path:        path,
		Name:        filepath.Base(path),
		Ext:         filepath.Ext(path),
	}
}

func NewFileInfo(storageType StorageType, path string, size int64, isDir bool, modTime time.Time) *FileInfo {
	fi := NewBasicFileInfo(storageType, path)
	fi.Size = size
	fi.IsDir = isDir
	fi.ModTime = modTime
	return fi
}

func (fi *FileInfo) LowerExt() string {
	return strings.ToLower(fi.Ext)
}

func (fi *FileInfo) String() string {
	return fi.StorageType.String() + ":" + fi.Path
}
