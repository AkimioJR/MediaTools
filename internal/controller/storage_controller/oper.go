package storage_controller

import (
	"MediaTools/internal/schemas"
	"io"
)

func Exists(file *schemas.FileInfo) (bool, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return false, schemas.ErrStorageProviderNotFound
	}
	return provider.Exists(file.Path)
}

func Mkdir(file *schemas.FileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return schemas.ErrStorageProviderNotFound
	}
	return provider.Mkdir(file.Path)
}

func Delete(file *schemas.FileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return schemas.ErrStorageProviderNotFound
	}
	return provider.Delete(file.Path)
}

func CreateFile(file *schemas.FileInfo, reader io.Reader) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return schemas.ErrStorageProviderNotFound
	}
	return provider.CreateFile(file.Path, reader)
}

func ReadFile(file *schemas.FileInfo) (io.ReadCloser, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(file.StorageType)
	if !exists {
		return nil, schemas.ErrStorageProviderNotFound
	}
	return provider.ReadFile(file.Path)
}

func List(dir *schemas.FileInfo) ([]schemas.FileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(dir.StorageType)
	if !exists {
		return nil, schemas.ErrStorageProviderNotFound
	}
	return provider.List(dir.Path)
}

func Copy(srcFile *schemas.FileInfo, dstFile *schemas.FileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	srcProvider, exists := getStorageProvider(srcFile.StorageType)
	if !exists {
		return schemas.ErrStorageProviderNotFound
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

func Move(srcFile *schemas.FileInfo, dstFile *schemas.FileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	srcProvider, exists := getStorageProvider(srcFile.StorageType)
	if !exists {
		return schemas.ErrStorageProviderNotFound
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

func Link(srcFile *schemas.FileInfo, dstFile *schemas.FileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	if srcFile.StorageType != dstFile.StorageType {
		return schemas.ErrNoSupport
	}

	provider, exists := getStorageProvider(srcFile.StorageType)
	if !exists {
		return schemas.ErrStorageProviderNotFound
	}
	return provider.Link(srcFile.Path, dstFile.Path)
}

func SoftLink(srcFile *schemas.FileInfo, dstFile *schemas.FileInfo) error {
	lock.RLock()
	defer lock.RUnlock()

	provider, exists := getStorageProvider(srcFile.StorageType)
	if !exists {
		return schemas.ErrStorageProviderNotFound
	}
	return provider.SoftLink(srcFile.Path, dstFile.Path)
}
