package storage

import (
	"io"
	"iter"
)

type StorageProvider interface {
	Init(config map[string]string) error // 初始化文件系统
	GetType() StorageType                // 获取文件系统类型
	GetTransferType() []TransferType     // 获取支持的传输类型

	// 路径级操作
	GetDetail(path string) (*StorageFileInfo, error) // 获取文件或目录的详细信息
	Exist(path string) (bool, error)                 // 判断文件是否存在
	Mkdir(path string) error                         // 创建目录（如果父目录不存在也需要创建）
	Delete(path string) error                        // 删除文件或目录
	Rename(oldPath string, newName string) error     // 重命名文件或目录

	// 文件内容操作
	CreateFile(path string, reader io.Reader) error // 创建文件并写入内容（如果父目录不存在也需要创建）
	ReadFile(path string) (io.ReadCloser, error)    // 读取文件内容

	// 目录操作
	ListRoot() (iter.Seq2[string, error], error)        // 列出根目录下的所有文件
	List(path string) (iter.Seq2[string, error], error) // 列出目录下的所有文件

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
