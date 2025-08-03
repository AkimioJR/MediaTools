package library_controller

import (
	"MediaTools/extensions"
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"slices"

	"github.com/sirupsen/logrus"
)

// 整理一个视频文件及其相关的字幕和音轨文件到指定目录
// srcFile: 源文件
// dstDir: 目标目录
// transferType: 传输类型（复制、移动、链接等）
// item: 媒体项（包含元数据）
// info: 识别到的媒体信息（为nil代表无需刮削）
// 返回值: 目标文件信息和可能的错误
func ArchiveMedia(
	srcFile *schemas.FileInfo,
	dstDir *schemas.FileInfo,
	transferType schemas.TransferType,
	item *schemas.MediaItem,
	info *schemas.MediaInfo,
) (*schemas.FileInfo, error) {
	transferLock.Lock()
	defer transferLock.Unlock()

	targetName, err := item.Format()
	if err != nil {
		return nil, err
	}
	dstFile := storage_controller.Join(dstDir, targetName)
	logrus.Infof("开始转移媒体文件：%s -> %s，转移类型类型：%s", srcFile, dstFile, transferType)

	err = storage_controller.TransferFile(srcFile, dstFile, transferType)
	if err != nil {
		return nil, err
	}

	{
		logrus.Info("开始转移字幕/音轨文件")
		srcDir := storage_controller.GetParent(srcFile)
		fileInfos, err := storage_controller.List(srcDir)
		if err != nil {
			logrus.Warningf("读取目录失败，跳过转移字幕/音轨文件：%v", err)
		} else {
			exts := append(extensions.SubtitleExtensions, extensions.AudioTrackExtensions...)
			for _, fi := range fileInfos {
				if fi.IsDir || fi.Path == srcDir.Path {
					continue // 跳过目录和源文件本身
				}

				if slices.Contains(exts, fi.LowerExt()) {
					otherDstFilePath := utils.ChangeExt(dstFile.Path, fi.Ext)
					otherDstFile, err := storage_controller.GetFile(otherDstFilePath, dstFile.StorageType)
					if err != nil {
						logrus.Warningf("获取文件 %s:%s 失败: %v", dstFile.StorageType, otherDstFilePath, err)
						continue
					}
					logrus.Debugf("转移字幕/音轨文件：%s -> %s", fi.String(), otherDstFile)
					err = storage_controller.Copy(&fi, otherDstFile) // 复制字幕/音轨文件
					if err != nil {
						logrus.Warningf("转移字幕/音轨文件失败：%v", err)
					}
				}
			}

		}
	}

	if info != nil {
		logrus.Info("开始生成刮削元数据")
		err = scrape_controller.Scrape(dstFile, info)
		if err != nil {
			logrus.Warningf("刮削数据失败：%v", err)
		}
	}
	return dstFile, nil
}
