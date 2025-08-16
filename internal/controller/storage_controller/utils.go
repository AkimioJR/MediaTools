package storage_controller

import (
	"MediaTools/internal/errs"
	"MediaTools/internal/schemas/storage"
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
