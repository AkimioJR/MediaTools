package storage_controller

import (
	"MediaTools/internal/errs"
	"MediaTools/internal/schemas/storage"
	"path/filepath"
)

func GetFile(path string, storageType storage.StorageType) (*storage.StorageFileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	_, exists := getStorageProvider(storageType)
	if !exists {
		return nil, errs.ErrStorageProviderNotFound
	}
	fi := storage.NewBasicFileInfo(storageType, path)
	return fi, nil
}

func GetParent(file *storage.StorageFileInfo) *storage.StorageFileInfo {
	parentPath := filepath.Dir(file.Path)
	return storage.NewBasicFileInfo(file.StorageType, parentPath)
}

func Join(file *storage.StorageFileInfo, elem ...string) *storage.StorageFileInfo {
	paths := make([]string, len(elem)+1)
	paths = append(paths, file.Path)
	paths = append(paths, elem...)
	path := filepath.Join(paths...)
	return storage.NewBasicFileInfo(file.StorageType, path)
}
