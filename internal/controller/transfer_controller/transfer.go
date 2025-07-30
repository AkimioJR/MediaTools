package transfer_controller

import (
	"MediaTools/extensions"
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
	transferType schemas.TransferType,
	item *schemas.MediaItem,
	info *schemas.MediaInfo,
) error {
	transferLock.Lock()
	defer transferLock.Unlock()

	targetName, err := item.Format()
	if err != nil {
		return err
	}
	targetPath := filepath.Join(targetDir, targetName)
	logrus.Infof("开始转移媒体文件：%s -> %s，类型：%s", srcPath, targetPath, transferType)

	err = TransferFile(srcPath, targetPath, transferType)
	if err != nil {
		return err
	}

	{
		logrus.Info("开始转移字幕/音轨文件")
		srcDir := filepath.Dir(srcPath)
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			logrus.Warningf("读取目录失败，跳过转移字幕/音轨文件：%v", err)
		}
		exts := append(extensions.SubtitleExtensions, extensions.AudioTrackExtensions...)
		for _, entry := range entries {
			if entry.IsDir() || entry.Name() == filepath.Base(srcPath) {
				continue // 跳过目录和源文件本身
			}

			if slices.Contains(exts, strings.ToLower(filepath.Ext(entry.Name()))) {
				srcFilePath := filepath.Join(srcDir, entry.Name())
				dstFilePath := utils.ChangeExt(targetPath, filepath.Ext(entry.Name()))
				logrus.Debugf("转移字幕/音轨文件：%s -> %s", srcFilePath, dstFilePath)
				err = CopyFile(srcFilePath, dstFilePath) // 复制字幕/音轨文件
				if err != nil {
					logrus.Warningf("转移字幕/音轨文件失败：%v", err)
				}
			}
		}
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
