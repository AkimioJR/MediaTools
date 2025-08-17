package library_controller

import (
	"MediaTools/extensions"
	"MediaTools/internal/controller/recognize_controller"
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"MediaTools/utils"
	"fmt"
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
	srcFile storage.StoragePath,
	dstDir storage.StoragePath,
	transferType storage.TransferType,
	item *schemas.MediaItem,
	info *schemas.MediaInfo,
) (storage.StoragePath, error) {
	lock.RLock()
	defer lock.RUnlock()

	targetName, err := recognize_controller.FormatVideo(item)
	if err != nil {
		return nil, err
	}
	dstPath := dstDir.Join(targetName)
	logrus.Infof("开始转移媒体文件：%s -> %s，转移类型类型：%s", srcFile, dstPath, transferType)

	err = storage_controller.TransferFile(srcFile, dstPath, transferType)
	if err != nil {
		return nil, err
	}

	{
		logrus.Info("开始转移字幕/音轨文件")
		srcDir := srcFile.Parent()
		paths, err := storage_controller.List(srcDir)
		if err != nil {
			logrus.Warningf("读取目录失败，跳过转移字幕/音轨文件：%v", err)
		} else {
			exts := append(extensions.SubtitleExtensions, extensions.AudioTrackExtensions...)
			for path, err := range paths {
				if err != nil {
					logrus.Warningf("遍历目录 %s 失败，跳过转移字幕/音轨文件：%v", srcDir, err)
					continue // 跳过错误的路径
				}
				if path.GetPath() == srcDir.GetPath() {
					continue // 跳过源目录本身
				}
				info, err := storage_controller.GetDetail(path)
				if err != nil {
					logrus.Warningf("获取文件 %s 详情失败，跳过转移字幕/音轨文件：%v", path, err)
					continue // 跳过获取详情失败的文件
				}
				if info.FileType == storage.FileTypeDirectory {
					logrus.Debugf("跳过目录：%s", info.Path)
					continue // 跳过目录
				}

				if slices.Contains(exts, info.LowerExt()) {
					otherdstPathPath := utils.ChangeExt(dstPath.GetPath(), info.Ext)
					otherdstPath, err := storage_controller.GetPath(otherdstPathPath, dstPath.GetStorageType())
					if err != nil {
						logrus.Warningf("获取文件 %s:%s 失败: %v", dstPath.GetStorageType(), otherdstPathPath, err)
						continue
					}
					logrus.Debugf("转移字幕/音轨文件：%s -> %s", info.String(), otherdstPath)
					err = storage_controller.TransferFile(info, otherdstPath, transferType) // 转移字幕或音轨文件
					if err != nil {
						logrus.Warningf("转移字幕/音轨文件失败：%v", err)
					}
				}
			}

		}
	}

	if info != nil {
		logrus.Info("开始生成刮削元数据")
		dstFile, err := storage_controller.GetDetail(dstPath)
		if err != nil {
			return nil, fmt.Errorf("获取目标文件路径失败：%w", err)
		}
		err = scrape_controller.Scrape(dstFile, info)
		if err != nil {
			logrus.Warningf("刮削数据失败：%v", err)
		}
	}
	return dstPath, nil
}

// ArchiveMediaSmart 处理媒体文件，智能识别并归档
// src: 源文件或目录
// 返回值: 可能的错误
// 注意：如果是目录，会递归处理目录下的所有文件
func ArchiveMediaSmart(src *storage.StorageFileInfo) error {
	lock.RLock()
	defer lock.RUnlock()
	var successNum int
	fn := func(file *storage.StorageFileInfo) error {
		libConfig := MatchLibrary(file)
		if libConfig == nil {
			return fmt.Errorf("未找到媒体库配置，跳过文件：%s", file.String())
		}

		logrus.Info("正在解析视频元数据：", src.Name)
		videoMeta, _, _ := recognize_controller.ParseVideoMeta(src.Name)
		info, err := tmdb_controller.RecognizeAndEnrichMedia(videoMeta)
		if err != nil {
			return fmt.Errorf("识别媒体信息失败：%w", err)
		}

		libraryBaseDir, err := storage_controller.GetPath(libConfig.DstPath, libConfig.DstType)
		if err != nil {
			return err
		}
		dstDir := libraryBaseDir.Join(GenFloder(libConfig, info)...)

		item, err := schemas.NewMediaItem(videoMeta, info)
		if err != nil {
			return fmt.Errorf("创建媒体项失败：%w", err)
		}
		_, err = ArchiveMedia(src, dstDir, libConfig.TransferType, item, info)
		if err != nil {
			return fmt.Errorf("转移媒体文件失败：%w", err)
		}
		successNum++
		return nil
	}

	if src.FileType == storage.FileTypeDirectory {
		logrus.Info("正在处理目录：", src.Path)
		iter, err := storage_controller.IterFiles(src)
		if err != nil {
			return fmt.Errorf("遍历目录 %s 失败：%w", src.Path, err)
		}
		for file, err := range iter {
			if err != nil {
				logrus.Warningf("处理文件 %s 时出错：%v", file.Path, err)
				continue
			}
			if file.FileType == storage.FileTypeDirectory {
				logrus.Debugf("跳过目录：%s", file.Path)
				continue
			}
			err = fn(file)
			if err != nil {
				logrus.Warningf("处理文件 %s 时出错：%v", file.Path, err)
			}
		}
		logrus.Infof("目录「%s」处理完成，成功：%d，", src.String(), successNum)
	} else {
		logrus.Info("正在处理文件：", src.Path)
		err := fn(src)
		if err != nil {
			return fmt.Errorf("处理文件 %s 时出错：%w", src.Path, err)
		}
		logrus.Info("媒体文件处理完成：", src.Path)
	}

	return nil
}
