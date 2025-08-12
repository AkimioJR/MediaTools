package storage_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/storage/local"
	"MediaTools/internal/schemas"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	lock             = sync.RWMutex{}
	storageProviders = make(map[schemas.StorageType]schemas.StorageProvider)
)

func Init() error {
	logrus.Info("开始初始化 Storage Controller...")
	for _, storageConfig := range config.Storages {
		_, err := RegisterStorageProvider(storageConfig)
		if err != nil {
			if storageConfig.Type != schemas.StorageUnknown {
				return fmt.Errorf("初始化 %s 存储器失败: %v", storageConfig.Type, err)
			} else {
				return err
			}
		}
	}

	logrus.Info("Storage Controller 初始化完成")
	return nil
}

func RegisterStorageProvider(c config.StorageConfig) (*schemas.StorageProviderItem, error) {
	lock.Lock()
	defer lock.Unlock()

	logrus.Debugf("开始初始化 %s 存储器...", c.Type)
	var provider schemas.StorageProvider
	switch c.Type {
	case schemas.StorageLocal:
		provider = &local.LocalStorage{}
		storageProviders[schemas.StorageLocal] = provider
	case schemas.StorageUnknown:
		return nil, fmt.Errorf("未知的存储类型: %s", c.Type)
	}
	err := provider.Init(c.Data)
	if err != nil {
		return nil, err
	}
	logrus.Infof("%s 存储器已注册", c.Type)
	item := schemas.NewStorageProviderItem(provider)
	return &item, nil
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

func GetStorageProvider(storageType schemas.StorageType) (*schemas.StorageProviderItem, error) {
	provider, exists := getStorageProvider(storageType)
	if !exists {
		return nil, fmt.Errorf("存储器 %s 不存在", storageType)
	}
	item := schemas.NewStorageProviderItem(provider)
	return &item, nil
}

func UnRegisterStorageProvider(storageType schemas.StorageType) (*schemas.StorageProviderItem, error) {
	lock.Lock()
	defer lock.Unlock()

	provider, exists := storageProviders[storageType]
	if !exists {
		return nil, fmt.Errorf("存储器 %s 不存在", storageType)
	}

	item := schemas.NewStorageProviderItem(provider)
	delete(storageProviders, storageType)
	logrus.Infof("已删除存储器: %s", storageType)
	return &item, nil
}
