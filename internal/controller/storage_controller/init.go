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
	lock.Lock()
	defer lock.Unlock()

	logrus.Info("开始初始化 Storage Controller...")
	for _, storageConfig := range config.Storages {
		err := RegisterStorageProvider(storageConfig)
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

func RegisterStorageProvider(c config.StorageConfig) error {
	logrus.Debugf("开始初始化 %s 存储器...", c.Type)
	switch c.Type {
	case schemas.StorageLocal:
		localStorage := &local.LocalStorage{}
		err := localStorage.Init(c.Data)
		if err != nil {
			return err
		}
		storageProviders[schemas.StorageLocal] = localStorage
	default:
		return fmt.Errorf("不支持的存储类型: %s", c.Type)
	}
	logrus.Infof("%s 存储器已注册", c.Type)
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
