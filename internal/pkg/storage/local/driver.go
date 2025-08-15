package local

import (
	"MediaTools/internal/errs"
	"MediaTools/internal/schemas/storage"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
}

func (s *LocalStorage) Init(config map[string]string) error {
	return nil
}

func (s *LocalStorage) GetType() storage.StorageType {
	return storage.StorageLocal
}

func (s *LocalStorage) GetTransferType() []storage.TransferType {
	return []storage.TransferType{storage.TransferCopy, storage.TransferMove, storage.TransferLink, storage.TransferSoftLink}
}

func (*LocalStorage) GetDetail(path string) (*storage.StorageFileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errs.ErrFileNotFound
		}
		return nil, err
	}

	return storage.NewFileInfo(
		storage.StorageLocal,
		path,
		info.Size(),
		info.IsDir(),
		info.ModTime(),
	), nil
}

func (s *LocalStorage) Exist(path string) (bool, error) {
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

func (s *LocalStorage) Rename(oldPath string, newName string) error {
	return errs.ErrStorageProvideNoSupport
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

func (s *LocalStorage) List(path string) ([]storage.StorageFileInfo, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileInfos []storage.StorageFileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		filePath := filepath.Join(path, file.Name())
		fileInfos = append(fileInfos, *storage.NewFileInfo(
			storage.StorageLocal,
			filePath,
			info.Size(),
			info.IsDir(),
			info.ModTime(),
		))
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

var _ storage.StorageProvider = (*LocalStorage)(nil)
