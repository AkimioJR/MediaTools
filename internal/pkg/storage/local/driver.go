package local

import (
	"MediaTools/internal/schemas"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
}

func (s *LocalStorage) Init(config map[string]any) error {
	return nil
}

func (s *LocalStorage) GetType() schemas.StorageType {
	return schemas.StorageLocal
}

func (s *LocalStorage) GetTransferType() []schemas.TransferType {
	return []schemas.TransferType{schemas.TransferCopy, schemas.TransferMove, schemas.TransferLink, schemas.TransferSoftLink}
}

func (s *LocalStorage) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *LocalStorage) Mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func (s *LocalStorage) Delete(path string) error {
	return os.RemoveAll(path)
}

func (s *LocalStorage) CreateFile(path string, reader io.Reader) error {
	err := s.Mkdir(filepath.Dir(path))
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	return err
}

func (s *LocalStorage) ReadFile(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *LocalStorage) List(path string) ([]schemas.FileInfo, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileInfos []schemas.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		fileInfos = append(fileInfos, schemas.FileInfo{
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime(),
		})
	}
	return fileInfos, nil
}

func (s *LocalStorage) Copy(srcPath string, dstPath string) error {
	err := s.Mkdir(filepath.Dir(dstPath))
	if err != nil {
		return err
	}
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	return s.CreateFile(dstPath, srcFile)
}

func (s *LocalStorage) Move(srcPath string, dstPath string) error {
	err := s.Mkdir(filepath.Dir(dstPath))
	if err != nil {
		return err
	}
	return os.Rename(srcPath, dstPath)
}

func (s *LocalStorage) Link(srcPath string, dstPath string) error {
	err := s.Mkdir(filepath.Dir(dstPath))
	if err != nil {
		return err
	}
	return os.Link(srcPath, dstPath)
}

func (s *LocalStorage) SoftLink(srcPath string, dstPath string) error {
	err := s.Mkdir(filepath.Dir(dstPath))
	if err != nil {
		return err
	}
	return os.Symlink(srcPath, dstPath)
}

var _ schemas.StorageProvider = (*LocalStorage)(nil)
