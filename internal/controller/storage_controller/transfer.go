package storage_controller

import (
	"MediaTools/internal/schemas"
	"fmt"
)

func TransferFile(srcFile *schemas.FileInfo, dstFile *schemas.FileInfo, transferType schemas.TransferType) error {
	if srcFile.StorageType != dstFile.StorageType &&
		(transferType == schemas.TransferLink || transferType == schemas.TransferSoftLink) {
		return fmt.Errorf("不支持使用转移方式 %s 将 %s 转移到 %s", transferType, srcFile, dstFile)
	}

	var err error
	switch transferType {
	case schemas.TransferCopy:
		err = Copy(srcFile, dstFile)
	case schemas.TransferMove:
		err = Move(srcFile, dstFile)
	case schemas.TransferLink:
		err = Link(srcFile, dstFile)
	case schemas.TransferSoftLink:
		err = SoftLink(srcFile, dstFile)
	default:
		err = fmt.Errorf("未知传输方式")
	}
	if err != nil {
		return fmt.Errorf("使用转移方式 %s 将 %s 转移到 %s 失败: %v", transferType, srcFile, dstFile, err)
	}
	return nil
}
