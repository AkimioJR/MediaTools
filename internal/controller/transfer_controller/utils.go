package transfer_controller

import (
	"MediaTools/internal/pkg/storage/model"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}

func TransferFile(srcPath string, targetPath string, transferType model.TransferType) error {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	var err error
	switch transferType {
	case model.TransferCopy:
		err = CopyFile(srcPath, targetPath)
	case model.TransferMove:
		err = os.Rename(srcPath, targetPath)
	case model.TransferLink:
		err = os.Link(srcPath, targetPath)
	case model.TransferSoftLink:
		err = os.Symlink(srcPath, targetPath)
	}
	return err
}
