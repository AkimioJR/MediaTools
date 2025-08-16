//go:build windows
// +build windows

package local

import (
	"fmt"
	"iter"

	"golang.org/x/sys/windows"
)

// ListRoot 列出根目录下的所有文件和目录
// 适用于 Windows 系统的根目录
func (s *LocalStorage) ListRoot() (iter.Seq2[string, error], error) {
	drives, err := windows.GetLogicalDrives()
	if err != nil {
		return nil, fmt.Errorf("failed to get logical drives: %w", err)
	}

	if drives == 0 {
		return nil, fmt.Errorf("no drives found")
	}

	return func(yield func(string, error) bool) {
		for i := range 26 { // 遍历所有驱动器
			if drives&(1<<uint(i)) != 0 {
				if !yield(fmt.Sprintf("%c:/", 'A'+i), nil) {
					return
				}
			}
		}
	}, nil
}
