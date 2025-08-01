package storage_controller

import (
	"MediaTools/internal/pkg/storage/local"
	"MediaTools/internal/schemas"
	"sync"
)

var (
	lock             = sync.RWMutex{}
	storageProviders = make(map[schemas.StorageType]schemas.StorageProvider)
)

func Init() error {
	lock.Lock()
	defer lock.Unlock()

	// 注册本地存储提供者
	localStorage := &local.LocalStorage{}
	err := localStorage.Init(nil)
	if err != nil {
		return err
	}
	storageProviders[schemas.StorageLocal] = localStorage

	return nil
}

func getStorageProvider(storageType schemas.StorageType) (schemas.StorageProvider, bool) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := storageProviders[storageType]
	return provider, exists
}

func ListStorageProviders() []schemas.StorageProviderItem {
	lock.RLock()
	defer lock.RUnlock()

	providers := make([]schemas.StorageProviderItem, 0, len(storageProviders))
	for _, provider := range storageProviders {
		providers = append(providers, schemas.NewStorageProviderItem(provider))
	}
	return providers
}
