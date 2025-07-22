package local

import (
	"MediaTools/internal/pkg/filesystem/model"
	"io"
	"os"
	pathlib "path"
)

type FileObj struct {
	path  string
	size  int64
	isDir bool
}

func (f *FileObj) GetType() model.StorageType {
	return model.StorageLocal
}

func (f *FileObj) GetName() string {
	return pathlib.Base(f.path)
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

func (f *FileObj) ReadContent() (reader io.ReadCloser, err error) {
	return os.Open(f.path)
}

var _ model.FileObject = (*FileObj)(nil)
