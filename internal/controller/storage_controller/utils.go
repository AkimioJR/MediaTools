package storage_controller

import (
	"MediaTools/internal/schemas"
	"path/filepath"
)

func GetFile(path string, storageType schemas.StorageType) (*schemas.FileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	_, exists := getStorageProvider(storageType)
	if !exists {
		return nil, schemas.ErrStorageProviderNotFound
	}
	fi := schemas.NewBasicFileInfo(storageType, path)
	return fi, nil
}

func GetParent(file *schemas.FileInfo) *schemas.FileInfo {
	parentPath := filepath.Dir(file.Path)
	return schemas.NewBasicFileInfo(file.StorageType, parentPath)
}

func Join(file *schemas.FileInfo, elem ...string) *schemas.FileInfo {
	paths := make([]string, len(elem)+1)
	paths = append(paths, file.Path)
	paths = append(paths, elem...)
	path := filepath.Join(paths...)
	return schemas.NewBasicFileInfo(file.StorageType, path)
}
