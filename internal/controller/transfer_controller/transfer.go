package transfer_controller

import (
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/pkg/storage/model"
	"MediaTools/internal/schemas"
	"log"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// 识别并整理一个文件
// srcPath: 源文件路径
// targetDir: 目标目录
// transferType: 传输类型（复制、移动、链接等）
// item: 媒体项（包含元数据）
// info: 识别到的媒体信息（为nil代表无需刮削）
func TransferMedia(
	srcPath string,
	targetDir string,
	transferType model.TransferType,
	item *schemas.MediaItem,
	info *schemas.MediaInfo,
) error {
	transferLock.Lock()
	defer transferLock.Unlock()

	log.Printf("开始转移媒体文件：%s 到 %s，类型：%s", srcPath, targetDir, transferType)

	targetName, err := item.Format()
	if err != nil {
		return err
	}
	targetPath := filepath.Join(targetDir, targetName)

	err = TransferFile(srcPath, targetPath, transferType)
	if err != nil {
		return err
	}

	if info != nil {
		logrus.Info("开始生成刮削元数据")
		err = scrape_controller.Scrape(targetPath, info)
		if err != nil {
			logrus.Warningf("刮削数据失败：%v", err)
		}
	}
	return nil
}
