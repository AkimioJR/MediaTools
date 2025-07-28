package local

import (
	"MediaTools/internal/pkg/storage/model"
)

type FileObj struct {
	name  string
	path  string
	size  int64
	isDir bool
}

func (f *FileObj) GetType() model.StorageType {
	return model.StorageLocal
}

func (f *FileObj) GetName() string {
	return f.name
}

func (f *FileObj) GetPath() string {
	return f.path
}

func (f *FileObj) IsDir() bool {
	return f.isDir
}

func (f *FileObj) GetSize() int64 {
	return f.size
}

var _ model.FileObject = (*FileObj)(nil)
