package model

import (
	"io"
	"strings"
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
		return "unknown"
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

type StorageBackend interface {
	Init(config map[string]any) error // 初始化文件系统
	GetType() StorageType             // 获取文件系统类型
	GetTransferType() []TransferType  // 获取支持的传输类型

	List(obj FileObject) ([]FileObject, error)         // 获取目录下的文件列表
	GetFile(path string) FileObject                    // 获取文件对象句柄
	CreateFile(obj FileObject, reader io.Reader) error // 创建文件
	ReadFile(obj FileObject) (io.ReadCloser, error)    // 读取文件内容
	Delete(obj FileObject) error                       // 删除文件或目录
	Rename(obj FileObject, newName string) error       // 重命名文件或目录
	Exists(obj FileObject) (bool, error)               // 检查文件或目录是否存在
	Mkdir(obj FileObject) error                        // 创建目录

	Copy(src FileObject, dst FileObject, storage StorageBackend) error // 复制文件或目录
	Move(src FileObject, dst FileObject, storage StorageBackend) error // 移动文件或目录
	Link(src FileObject, dst FileObject) error                         // 创建硬链接
	SoftLink(src FileObject, dst FileObject) error                     // 创建软链接
}
