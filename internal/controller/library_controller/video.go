package library_controller

import (
	"MediaTools/extensions"
	"MediaTools/internal/controller/format_controller"
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
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
	srcFile *schemas.FileInfo,
	dstDir *schemas.FileInfo,
	transferType schemas.TransferType,
	item *schemas.MediaItem,
	info *schemas.MediaInfo,
) (*schemas.FileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	targetName, err := format_controller.FormatVideo(item)
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
					err = storage_controller.TransferFile(&fi, otherDstFile, transferType) // 转移字幕或音轨文件
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

// ArchiveMediaSmart 处理媒体文件，智能识别并归档
// src: 源文件或目录
// 返回值: 可能的错误
// 注意：如果是目录，会递归处理目录下的所有文件
func ArchiveMediaSmart(src *schemas.FileInfo) error {
	lock.RLock()
	defer lock.RUnlock()
	var successNum int
	fn := func(file *schemas.FileInfo) error {
		libConfig := MatchLibrary(file)
		if libConfig == nil {
			return fmt.Errorf("未找到媒体库配置，跳过文件：%s", file.String())
		}

		logrus.Info("正在解析视频元数据：", src.Name)
		videoMeta := meta.ParseVideoMeta(src.Name)
		info, err := tmdb_controller.RecognizeAndEnrichMedia(videoMeta, nil, nil)
		if err != nil {
			return fmt.Errorf("识别媒体信息失败：%w", err)
		}

		libraryBaseDir, err := storage_controller.GetFile(libConfig.DstPath, libConfig.DstType)
		if err != nil {
			return err
		}
		dstDir := storage_controller.Join(libraryBaseDir, GenFloder(libConfig, info)...)

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

	if src.IsDir {
		logrus.Info("正在处理目录：", src.Path)
		err := storage_controller.IterFiles(src, fn)
		if err != nil {
			return fmt.Errorf("处理目录 %s 时出错，已成功：%d 个，错误：%w", src.Path, successNum, err)
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
