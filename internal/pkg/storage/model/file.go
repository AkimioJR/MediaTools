package model

import "time"

type FileInfo struct {
	StorageType StorageType
	Name        string
	Path        string
	Size        int64
	IsDir       bool
	ModTime     time.Time
}
