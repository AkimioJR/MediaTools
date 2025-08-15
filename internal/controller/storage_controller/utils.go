package storage_controller

import (
	"MediaTools/internal/errs"
	"MediaTools/internal/schemas/storage"
	"path/filepath"
)

func GetPath(path string, storageType storage.StorageType) (storage.StoragePath, error) {
	lock.RLock()
	defer lock.RUnlock()

	_, exists := getStorageProvider(storageType)
	if !exists {
		return nil, errs.ErrStorageProviderNotFound
	}
	return storage.NewStoragePath(storageType, path), nil
}

func GetParent(path storage.StoragePath) storage.StoragePath {
	parentPath := filepath.Dir(path.GetPath())
	return storage.NewStoragePath(path.GetStorageType(), parentPath)
}

func Join(file storage.StoragePath, elem ...string) storage.StoragePath {
	paths := make([]string, len(elem)+1)
	paths = append(paths, file.GetPath())
	paths = append(paths, elem...)
	return storage.NewStoragePath(file.GetStorageType(), filepath.Join(paths...))
}
