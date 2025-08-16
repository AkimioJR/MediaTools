package storage

import (
	pathlib "path"
	"strings"
	"time"
)

type StoragePath interface {
	GetStorageType() StorageType
	GetPath() string
	GetName() string
	GetExt() string
	LowerExt() string
	String() string

	Parent() StoragePath
	Join(elem ...string) StoragePath
}

func NewStoragePath(storageType StorageType, path string) StoragePath {
	return &StorageFileInfo{
		StorageType: storageType,
		Path:        pathlib.Clean(path),
		Name:        pathlib.Base(path),
		Ext:         pathlib.Ext(path),
	}
}

type StorageFileInfo struct {
	StorageType StorageType `json:"storage_type"` // 存储系统类型
	Path        string      `json:"path"`         // 文件路径
	Name        string      `json:"name"`         // 文件名
	Ext         string      `json:"ext"`          // 文件扩展名
	Size        int64       `json:"size"`         // 文件大小
	IsDir       bool        `json:"is_dir"`       // 是否为目录
	ModTime     time.Time   `json:"mod_time"`     // 文件修改时间
}

func NewFileInfo(storageType StorageType, path string, size int64, isDir bool, modTime time.Time) *StorageFileInfo {
	fi := NewStoragePath(storageType, path).(*StorageFileInfo)
	fi.Size = size
	fi.IsDir = isDir
	fi.ModTime = modTime
	return fi
}

func (fi *StorageFileInfo) Parent() StoragePath {
	return NewStoragePath(fi.StorageType, pathlib.Dir(fi.Path))
}

func (fi *StorageFileInfo) Join(elem ...string) StoragePath {
	paths := make([]string, len(elem)+1)
	paths[0] = fi.Path
	paths = append(paths, elem...)
	return NewStoragePath(fi.StorageType, pathlib.Join(paths...))
}

func (fi *StorageFileInfo) GetStorageType() StorageType {
	return fi.StorageType
}

func (fi *StorageFileInfo) GetPath() string {
	return fi.Path
}

func (fi *StorageFileInfo) GetName() string {
	return fi.Name
}

func (fi *StorageFileInfo) GetExt() string {
	return fi.Ext
}

func (fi *StorageFileInfo) LowerExt() string {
	return strings.ToLower(fi.Ext)
}

func (fi *StorageFileInfo) String() string {
	return fi.StorageType.String() + ":" + fi.Path
}

var _ StoragePath = (*StorageFileInfo)(nil)
