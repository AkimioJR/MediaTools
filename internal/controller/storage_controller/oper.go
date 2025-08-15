package storage_controller

import (
	"MediaTools/internal/errs"
	"MediaTools/internal/schemas/storage"
	"io"
	"path/filepath"
)

func GetDetail(path storage.StoragePath) (*storage.StorageFileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(path.GetStorageType())
	if !exists {
		return nil, errs.ErrStorageProviderNotFound
	}
	return provider.GetDetail(path.GetPath())
}

func Exist(path storage.StoragePath) (bool, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(path.GetStorageType())
	if !exists {
		return false, errs.ErrStorageProviderNotFound
	}
	return provider.Exist(path.GetPath())
}

func Mkdir(path storage.StoragePath) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(path.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.Mkdir(path.GetPath())
}

func Rename(path storage.StoragePath, newName string) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(path.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	err := provider.Rename(path.GetPath(), newName)
	switch err {
	case nil:
		return nil
	case errs.ErrStorageProvideNoSupport: // 如果不支持重命名，则尝试使用移动的方式
		return provider.Move(path.GetPath(), filepath.Join(filepath.Dir(path.GetPath()), newName))
	default:
		return err
	}
}

func Delete(path storage.StoragePath) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(path.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.Delete(path.GetPath())
}

func CreateFile(path storage.StoragePath, reader io.Reader) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(path.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.CreateFile(path.GetPath(), reader)
}

func ReadFile(path storage.StoragePath) (io.ReadCloser, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(path.GetStorageType())
	if !exists {
		return nil, errs.ErrStorageProviderNotFound
	}
	return provider.ReadFile(path.GetPath())
}

func List(dir storage.StoragePath) ([]storage.StorageFileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(dir.GetStorageType())
	if !exists {
		return nil, errs.ErrStorageProviderNotFound
	}
	return provider.List(dir.GetPath())
}

func Copy(srcPath storage.StoragePath, dstPath storage.StoragePath) error {
	lock.RLock()
	defer lock.RUnlock()

	srcProvider, exists := getStorageProvider(srcPath.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	if srcPath.GetStorageType() != dstPath.GetStorageType() {
		reader, err := srcProvider.ReadFile(srcPath.GetPath())
		if err != nil {
			return err
		}
		defer reader.Close()
		return CreateFile(dstPath, reader)
	}
	return srcProvider.Copy(srcPath.GetPath(), dstPath.GetPath())
}

func Move(srcPath storage.StoragePath, dstPath storage.StoragePath) error {
	lock.RLock()
	defer lock.RUnlock()

	srcProvider, exists := getStorageProvider(srcPath.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	if srcPath.GetStorageType() != dstPath.GetStorageType() {
		reader, err := srcProvider.ReadFile(srcPath.GetPath())
		if err != nil {
			return err
		}
		defer reader.Close()
		err = CreateFile(dstPath, reader)
		if err != nil {
			return err
		}
		return srcProvider.Delete(srcPath.GetPath())
	}
	return srcProvider.Move(srcPath.GetPath(), dstPath.GetPath())
}

func Link(srcPath storage.StoragePath, dstPath storage.StoragePath) error {
	lock.RLock()
	defer lock.RUnlock()

	if srcPath.GetStorageType() != dstPath.GetStorageType() {
		return errs.ErrStorageProvideNoSupport
	}

	provider, exists := getStorageProvider(srcPath.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.Link(srcPath.GetPath(), dstPath.GetPath())
}

func SoftLink(srcPath storage.StoragePath, dstPath storage.StoragePath) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(srcPath.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.SoftLink(srcPath.GetPath(), dstPath.GetPath())
}

func IterFiles(dir storage.StoragePath, fn func(file *storage.StorageFileInfo) error) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(dir.GetStorageType())
	if !exists {
		return errs.ErrStorageProviderNotFound
	}

	files, err := provider.List(dir.GetPath())
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir {
			IterFiles(&file, fn)
		}
		err := fn(&file)
		if err != nil {
			return err
		}
	}
	return nil
}
