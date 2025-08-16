package storage

import (
	"encoding/json"
	"strings"
)

type FileType uint8

const (
	FileTypeUnknown   FileType = iota // 未知类型
	FileTypeFile                      // 普通文件
	FileTypeDirectory                 // 目录
)

func (ft FileType) String() string {
	switch ft {
	case FileTypeFile:
		return "File"
	case FileTypeDirectory:
		return "Directory"
	default: // FileTypeUnknown
		return "UnknownFileType"
	}
}

func (ft *FileType) ParseString(s string) error {
	switch strings.ToLower(s) {
	case "file":
		*ft = FileTypeFile
	case "directory":
		*ft = FileTypeDirectory
	default:
		*ft = FileTypeUnknown
	}
	return nil
}

func (ft FileType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ft.String() + `"`), nil
}

func (ft *FileType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return ft.ParseString(s)
}
