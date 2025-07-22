package local

import (
	"MediaTools/internal/pkg/filesystem/model"
	"io"
	"os"
	pathlib "path"
)

type LocalStorage struct {
}

func (s *LocalStorage) Init(config map[string]any) error {
	return nil
}

func (s *LocalStorage) GetType() model.StorageType {
	return model.StorageLocal
}
func (s *LocalStorage) GetTransferType() []model.TransferType {
	return []model.TransferType{model.TransferCopy, model.TransferMove, model.TransferLink, model.TransferSoftLink}
}

func (s *LocalStorage) GetRoot() (model.FileObject, error) {
	info, err := os.Stat("/")
	if err != nil {
		return nil, err
	}
	obj := FileObj{
		path:  "/",
		size:  info.Size(),
		isDir: info.IsDir(),
	}
	return &obj, nil
}

func (s *LocalStorage) List(obj model.FileObject) ([]model.FileObject, error) {
	file, ok := obj.(*FileObj)
	if !ok {
		return nil, model.ErrNoSupport
	}

	entries, err := os.ReadDir(file.GetPath())
	if err != nil {
		return nil, err
	}
	var fileObjects []model.FileObject
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fileObjects = append(fileObjects, &FileObj{
			path:  pathlib.Join(file.GetPath(), info.Name()),
			size:  info.Size(),
			isDir: info.IsDir(),
		})
	}
	return fileObjects, nil
}

func (s *LocalStorage) NewFile(dir model.FileObject, name string) model.FileObject {
	fileObj := &FileObj{
		path:  pathlib.Join(dir.GetPath(), name),
		isDir: false,
	}
	return fileObj
}

func (s *LocalStorage) CreateFile(obj model.FileObject, reader io.Reader) error {
	fileObj, ok := obj.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}

	file, err := os.Create(fileObj.GetPath())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (s *LocalStorage) Delete(obj model.FileObject) error {
	fileObj, ok := obj.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}

	if fileObj.IsDir() {
		return os.RemoveAll(fileObj.GetPath())
	}
	return os.Remove(fileObj.GetPath())
}
func (s *LocalStorage) Rename(obj model.FileObject, newName string) error {
	fileObj, ok := obj.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}

	newPath := pathlib.Join(pathlib.Dir(fileObj.GetPath()), newName)
	return os.Rename(fileObj.GetPath(), newPath)
}

func (s *LocalStorage) Exists(obj model.FileObject) (bool, error) {
	fileObj, ok := obj.(*FileObj)
	if !ok {
		return false, model.ErrNoSupport
	}

	_, err := os.Stat(fileObj.GetPath())
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *LocalStorage) Mkdir(obj model.FileObject) error {
	fileObj, ok := obj.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}

	return os.MkdirAll(fileObj.GetPath(), os.ModePerm)
}

func (s *LocalStorage) Copy(src model.FileObject, dst model.FileObject, dstFS model.FileSystem) error {
	srcFile, ok := src.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}
	reader, err := srcFile.ReadContent()
	if err != nil {
		return err
	}
	defer reader.Close()
	return dstFS.CreateFile(dst, reader)
}

func (s *LocalStorage) Move(src model.FileObject, dst model.FileObject, dstFS model.FileSystem) error {
	_, ok1 := src.(*FileObj)
	_, ok2 := dst.(*FileObj)
	if ok1 && ok2 {
		srcFile, err := os.Open(src.GetPath())
		if err != nil {
			return err
		}
		defer srcFile.Close()
		dstFile, err := os.Create(dst.GetPath())
		if err != nil {
			return err
		}
		defer dstFile.Close()
		_, err = io.Copy(dstFile, srcFile)
		return err
	}
	if err := s.Copy(src, dst, dstFS); err != nil {
		return err
	}
	return s.Delete(src)
}

func (s *LocalStorage) Link(src model.FileObject, dst model.FileObject) error {
	srcFile, ok := src.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}
	dstFile, ok := dst.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}
	// Link 将新名称（newname）创建为旧名称（oldname）文件的硬链接。如果出现错误，该错误将为 *LinkError 类型。
	return os.Link(srcFile.GetPath(), dstFile.GetPath())
}

func (s *LocalStorage) SoftLink(src model.FileObject, dst model.FileObject) error {
	srcFile, ok := src.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}
	dstFile, ok := dst.(*FileObj)
	if !ok {
		return model.ErrNoSupport
	}

	// Symlink 将新名称（newname）创建为指向旧名称（oldname）的符号链接。
	// 在 Windows 系统上，若指向的旧名称（oldname）不存在，创建的符号链接会是文件类型的符号链接；
	// 即便之后将旧名称（oldname）创建为目录，该符号链接也无法正常工作。
	// 如果出现错误，错误类型将为 *LinkError。
	return os.Symlink(srcFile.GetPath(), dstFile.GetPath())
}

var _ model.FileSystem = (*LocalStorage)(nil)
