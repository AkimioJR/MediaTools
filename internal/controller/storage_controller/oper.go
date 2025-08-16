package storage_controller

import (
	"MediaTools/internal/errs"
	"MediaTools/internal/schemas/storage"
	"io"
	"iter"
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
		return provider.Move(path.GetPath(), path.Parent().Join(newName).GetPath())
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

func List(dir storage.StoragePath) (iter.Seq2[storage.StoragePath, error], error) {
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

func IterFiles(dir storage.StoragePath) (iter.Seq2[*storage.StorageFileInfo, error], error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(dir.GetStorageType())
	if !exists {
		return nil, errs.ErrStorageProviderNotFound
	}

	// 验证目录是否存在
	dirExists, err := provider.Exist(dir.GetPath())
	if err != nil {
		return nil, err
	}
	if !dirExists {
		return nil, errs.ErrStorageProviderNotFound
	}

	// 验证路径是否为目录
	dirInfo, err := provider.GetDetail(dir.GetPath())
	if err != nil {
		return nil, err
	}
	if !dirInfo.IsDir {
		return nil, errs.ErrStorageProviderNotFound // 或者定义一个新的错误类型
	}

	return func(yield func(*storage.StorageFileInfo, error) bool) {
		iterFilesRecursive(provider, dir.GetPath(), yield)
	}, nil
}

// iterFilesRecursive 递归遍历目录中的所有文件
func iterFilesRecursive(provider storage.StorageProvider, dirPath string, yield func(*storage.StorageFileInfo, error) bool) {
	iter, err := provider.List(dirPath)
	if err != nil {
		if !yield(nil, err) {
			return
		}
		return
	}

	for path, err := range iter {
		if err != nil {
			if !yield(nil, err) { // 如果迭代器被中断，则退出
				return
			}
			continue
		}

		info, err := provider.GetDetail(path.GetPath())
		if err != nil {
			if !yield(nil, err) {
				return
			}
			continue
		}

		if info.IsDir {
			// 如果是目录，递归遍历
			iterFilesRecursive(provider, info.Path, yield)
		} else {
			// 如果是文件，yield 返回
			if !yield(info, nil) {
				return
			}
		}
	}
}
