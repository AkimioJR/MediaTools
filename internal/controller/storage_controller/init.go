package storage_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/storage/local"
	"MediaTools/internal/schemas/storage"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	lock             = sync.RWMutex{}
	storageProviders = make(map[storage.StorageType]storage.StorageProvider)
)

func Init() error {
	logrus.Info("开始初始化 Storage Controller...")
	for _, storageConfig := range config.Storages {
		_, err := RegisterStorageProvider(storageConfig)
		if err != nil {
			if storageConfig.Type != storage.StorageUnknown {
				return fmt.Errorf("初始化 %s 存储器失败: %v", storageConfig.Type, err)
			} else {
				return err
			}
		}
	}

	logrus.Info("Storage Controller 初始化完成")
	return nil
}

func RegisterStorageProvider(c config.StorageConfig) (*storage.StorageProviderItem, error) {
	lock.Lock()
	defer lock.Unlock()

	logrus.Debugf("开始初始化 %s 存储器...", c.Type)
	var provider storage.StorageProvider
	switch c.Type {
	case storage.StorageLocal:
		provider = &local.LocalStorage{}
		storageProviders[storage.StorageLocal] = provider
	case storage.StorageUnknown:
		return nil, fmt.Errorf("未知的存储类型: %s", c.Type)
	}
	err := provider.Init(c.Data)
	if err != nil {
		return nil, err
	}
	logrus.Infof("%s 存储器已注册", c.Type)
	item := storage.NewStorageProviderItem(provider)
	return &item, nil
}

func getStorageProvider(storageType storage.StorageType) (storage.StorageProvider, bool) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := storageProviders[storageType]
	return provider, exists
}

func ListStorageProviders() []storage.StorageProviderItem {
	lock.RLock()
	defer lock.RUnlock()

	providers := make([]storage.StorageProviderItem, 0, len(storageProviders))
	for _, provider := range storageProviders {
		providers = append(providers, storage.NewStorageProviderItem(provider))
	}
	return providers
}

func GetStorageProvider(storageType storage.StorageType) (*storage.StorageProviderItem, error) {
	provider, exists := getStorageProvider(storageType)
	if !exists {
		return nil, fmt.Errorf("存储器 %s 不存在", storageType)
	}
	item := storage.NewStorageProviderItem(provider)
	return &item, nil
}

func UnRegisterStorageProvider(storageType storage.StorageType) (*storage.StorageProviderItem, error) {
	lock.Lock()
	defer lock.Unlock()

	provider, exists := storageProviders[storageType]
	if !exists {
		return nil, fmt.Errorf("存储器 %s 不存在", storageType)
	}

	item := storage.NewStorageProviderItem(provider)
	delete(storageProviders, storageType)
	logrus.Infof("已删除存储器: %s", storageType)
	return &item, nil
}
