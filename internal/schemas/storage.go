package schemas

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
	"time"
)

type StorageType uint8

const (
	StorageUnknown StorageType = iota // 未知文件系统
	StorageLocal                      // 本地文件系统
)

func (t StorageType) String() string {
	switch t {
	case StorageLocal:
		return "LocalStorage"
	default:
		return "UnknownStorage"
	}
}

func ParseStorageType(s string) StorageType {
	switch strings.ToLower(s) {
	case "localstorage":
		return StorageLocal
	default:
		return StorageUnknown
	}
}

func (t StorageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *StorageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = ParseStorageType(s)
	return nil
}

type TransferType uint8

const (
	TransferUnknown  TransferType = iota // 未知传输类型
	TransferCopy                         // 复制
	TransferMove                         // 移动
	TransferLink                         // 硬链接
	TransferSoftLink                     // 软链接
)

func (t TransferType) String() string {
	switch t {
	case TransferCopy:
		return "Copy"
	case TransferMove:
		return "Move"
	case TransferLink:
		return "Link"
	case TransferSoftLink:
		return "SoftLink"
	default:
		return "UnknownTransferType"
	}
}

func (t TransferType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TransferType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = ParseTransferType(s)
	return nil
}

func ParseTransferType(s string) TransferType {
	switch strings.ToLower(s) {
	case "copy":
		return TransferCopy
	case "move":
		return TransferMove
	case "link":
		return TransferLink
	case "softlink":
		return TransferSoftLink
	default:
		return TransferUnknown
	}
}

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

type StorageProvider interface {
	Init(config map[string]any) error // 初始化文件系统
	GetType() StorageType             // 获取文件系统类型
	GetTransferType() []TransferType  // 获取支持的传输类型

	// 路径级操作
	Exists(path string) (bool, error) // 判断文件是否存在
	Mkdir(path string) error          // 创建目录（如果父目录不存在也需要创建）
	Delete(path string) error         // 删除文件或目录

	// 文件内容操作
	CreateFile(path string, reader io.Reader) error // 创建文件并写入内容（如果父目录不存在也需要创建）
	ReadFile(path string) (io.ReadCloser, error)    // 读取文件内容

	// 目录操作
	List(path string) ([]FileInfo, error) // 列出目录下的所有文件

	// 文件传输操作
	Copy(srcPath string, dstPath string) error     // 复制文件
	Move(srcPath string, dstPath string) error     // 移动文件
	Link(srcPath string, dstPath string) error     // 硬链接文件
	SoftLink(srcPath string, dstPath string) error // 软链接文件
}

type StorageProviderItem struct {
	StorageType  StorageType    `json:"storage_type"`
	TransferType []TransferType `json:"transfer_type"`
}

func NewStorageProviderItem(provider StorageProvider) StorageProviderItem {
	return StorageProviderItem{
		StorageType:  provider.GetType(),
		TransferType: provider.GetTransferType(),
	}
}
