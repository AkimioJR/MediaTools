package storage_controller

import (
	"MediaTools/internal/schemas/storage"
	"fmt"
)

func TransferFile(srcPath storage.StoragePath, dstPath storage.StoragePath, transferType storage.TransferType) error {
	if srcPath.GetStorageType() != dstPath.GetStorageType() &&
		(transferType == storage.TransferLink || transferType == storage.TransferSoftLink) {
		return fmt.Errorf("不支持使用转移方式 %s 将 %s 转移到 %s", transferType, srcPath, dstPath)
	}

	var err error
	switch transferType {
	case storage.TransferCopy:
		err = Copy(srcPath, dstPath)
	case storage.TransferMove:
		err = Move(srcPath, dstPath)
	case storage.TransferLink:
		err = Link(srcPath, dstPath)
	case storage.TransferSoftLink:
		err = SoftLink(srcPath, dstPath)
	default:
		err = fmt.Errorf("未知传输方式")
	}
	if err != nil {
		return fmt.Errorf("使用转移方式 %s 将 %s 转移到 %s 失败: %v", transferType, srcPath, dstPath, err)
	}
	return nil
}
