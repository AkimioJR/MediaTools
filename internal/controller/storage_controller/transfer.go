package storage_controller

import (
	"MediaTools/internal/schemas/storage"
	"fmt"
)

func TransferFile(srcFile *storage.FileInfo, dstFile *storage.FileInfo, transferType storage.TransferType) error {
	if srcFile.StorageType != dstFile.StorageType &&
		(transferType == storage.TransferLink || transferType == storage.TransferSoftLink) {
		return fmt.Errorf("不支持使用转移方式 %s 将 %s 转移到 %s", transferType, srcFile, dstFile)
	}

	var err error
	switch transferType {
	case storage.TransferCopy:
		err = Copy(srcFile, dstFile)
	case storage.TransferMove:
		err = Move(srcFile, dstFile)
	case storage.TransferLink:
		err = Link(srcFile, dstFile)
	case storage.TransferSoftLink:
		err = SoftLink(srcFile, dstFile)
	default:
		err = fmt.Errorf("未知传输方式")
	}
	if err != nil {
		return fmt.Errorf("使用转移方式 %s 将 %s 转移到 %s 失败: %v", transferType, srcFile, dstFile, err)
	}
	return nil
}
