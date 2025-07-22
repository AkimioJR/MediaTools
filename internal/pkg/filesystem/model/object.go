package model

import "io"

type FileObject interface {
	GetType() StorageType                // 获取文件对象类型
	GetName() string                     // 获取文件名
	GetPath() string                     // 获取文件完整路径
	IsDir() bool                         // 是否为目录
	GetSize() int64                      // 获取文件大小
	ReadContent() (io.ReadCloser, error) // 读取文件内容
}
