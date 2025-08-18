//go:build !windows
// +build !windows

package local

import (
	"MediaTools/internal/schemas/storage"
	"iter"
)

// ListRoot 列出根目录下的所有文件和目录
// 适用于非 Windows 系统的根目录
func (s *LocalStorage) ListRoot() (iter.Seq2[storage.StorageEntry, error], error) {
	return s.List("/")
}
