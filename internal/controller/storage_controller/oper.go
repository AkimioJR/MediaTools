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

func Exists(path storage.StoragePath) (bool, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return false, errs.ErrStorageProviderNotFound
	}
	return provider.Exists(file.Path)
}

func Mkdir(file *storage.StorageFileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.Mkdir(file.Path)
}

func Rename(file *storage.StorageFileInfo, newName string) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	err := provider.Rename(file.Path, newName)
	switch err {
	case nil:
		return nil
	case errs.ErrStorageProvideNoSupport: // 如果不支持重命名，则尝试使用移动的方式
		return provider.Move(file.Path, filepath.Join(filepath.Dir(file.Path), newName))
	default:
		return err
	}
}

func Delete(file *storage.StorageFileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.Delete(file.Path)
}

func CreateFile(file *storage.StorageFileInfo, reader io.Reader) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.CreateFile(file.Path, reader)
}

func ReadFile(file *storage.StorageFileInfo) (io.ReadCloser, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return nil, errs.ErrStorageProviderNotFound
	}
	return provider.ReadFile(file.Path)
}

func List(dir *storage.StorageFileInfo) ([]storage.StorageFileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(dir.StorageType)
	if !exists {
		return nil, errs.ErrStorageProviderNotFound
	}
	return provider.List(dir.Path)
}

func Copy(srcFile *storage.StorageFileInfo, dstFile *storage.StorageFileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	srcProvider, exists := getStorageProvider(srcFile.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	if srcFile.StorageType != dstFile.StorageType {
		reader, err := srcProvider.ReadFile(srcFile.Path)
		if err != nil {
			return err
		}
		defer reader.Close()
		return CreateFile(dstFile, reader)
	}
	return srcProvider.Copy(srcFile.Path, dstFile.Path)
}

func Move(srcFile *storage.StorageFileInfo, dstFile *storage.StorageFileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	srcProvider, exists := getStorageProvider(srcFile.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	if srcFile.StorageType != dstFile.StorageType {
		reader, err := srcProvider.ReadFile(srcFile.Path)
		if err != nil {
			return err
		}
		defer reader.Close()
		err = CreateFile(dstFile, reader)
		if err != nil {
			return err
		}
		return srcProvider.Delete(srcFile.Path)
	}
	return srcProvider.Move(srcFile.Path, dstFile.Path)
}

func Link(srcFile *storage.StorageFileInfo, dstFile *storage.StorageFileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	if srcFile.StorageType != dstFile.StorageType {
		return errs.ErrStorageProvideNoSupport
	}

	provider, exists := getStorageProvider(srcFile.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.Link(srcFile.Path, dstFile.Path)
}

func SoftLink(srcFile *storage.StorageFileInfo, dstFile *storage.StorageFileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(srcFile.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}
	return provider.SoftLink(srcFile.Path, dstFile.Path)
}

func IterFiles(dir *storage.StorageFileInfo, fn func(file *storage.StorageFileInfo) error) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(dir.StorageType)
	if !exists {
		return errs.ErrStorageProviderNotFound
	}

	files, err := provider.List(dir.Path)
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
