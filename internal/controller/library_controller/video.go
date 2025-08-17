package library_controller

import (
	"MediaTools/extensions"
	"MediaTools/internal/controller/recognize_controller"
	"MediaTools/internal/controller/scrape_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/wordmatch"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"MediaTools/utils"
	"fmt"
	"slices"
	"strings"

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
				if info.Type == storage.FileTypeDirectory {
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

func ArchiveMediaAdvanced(srcFile storage.StoragePath, dstDir storage.StoragePath,
	transferType storage.TransferType, mediaType meta.MediaType,
	tmdbID int, season int, episodeStr string, episodeOffset string,
	part string, organizeByType bool, organizeByCategory bool, scrape bool,
) error {
	lock.RLock()
	defer lock.RUnlock()

	videoMeta, rule1, rule2 := recognize_controller.ParseVideoMeta(srcFile.GetName())
	switch {
	case rule1 != "" && rule2 != "":
		logrus.Debugf("解析视频元数据: %s，匹配的自定义规则：%s，应用的自定义媒体规则：%s", srcFile.GetName(), rule1, rule2)
	case rule1 != "":
		logrus.Debugf("解析视频元数据: %s，匹配的自定义规则：%s", srcFile.GetName(), rule1)
	case rule2 != "":
		logrus.Debugf("解析视频元数据: %s，应用的自定义媒体规则：%s", srcFile.GetName(), rule2)
	default:
		logrus.Debugf("解析视频元数据: %s，没有匹配到自定义规则和应用的自定义媒体规则", srcFile.GetName())
	}

	var msgs []string
	if mediaType != meta.MediaTypeUnknown {
		videoMeta.MediaType = mediaType
		msgs = append(msgs, fmt.Sprintf("媒体类型: %s", mediaType))
	}
	if tmdbID != 0 {
		videoMeta.TMDBID = tmdbID
		msgs = append(msgs, fmt.Sprintf("TMDB ID: %d", tmdbID))
	}
	if season >= -1 {
		videoMeta.Season = season
		msgs = append(msgs, fmt.Sprintf("季数: %d", season))
	}

	if episodeStr != "" {
		startEpisode, endEpisode, err := parseEpisodeStr(episodeStr)
		if err != nil {
			return fmt.Errorf("解析集数失败: %w", err)
		}
		videoMeta.Episode = startEpisode
		videoMeta.EndEpisode = endEpisode
		if endEpisode != 0 {
			msgs = append(msgs, fmt.Sprintf("集数范围: %d-%d", startEpisode, endEpisode))
		} else {
			msgs = append(msgs, fmt.Sprintf("集数: %d", startEpisode))
		}
	} else if episodeOffset != "" { // 当 episodeStr 为空时，才使用 episodeOffset
		offsetEpisode, err := wordmatch.ParseOffsetExpr(episodeOffset, videoMeta.Episode)
		if err != nil {
			return fmt.Errorf("解析集数偏移表达式失败：%w", err)
		}
		videoMeta.Episode = offsetEpisode
		msgs = append(msgs, fmt.Sprintf("集数偏移：%s，计算结果：%d", episodeOffset, offsetEpisode))
	}

	if part != "" {
		videoMeta.Part = part
		msgs = append(msgs, fmt.Sprintf("指定分段: %s", part))
	}
	if len(msgs) > 0 {
		logrus.Infof("更新 %s 媒体元数据：%s", srcFile.GetName(), strings.Join(msgs, ", "))
	}

	info, err := tmdb_controller.RecognizeAndEnrichMedia(videoMeta)
	if err != nil {
		return fmt.Errorf("识别媒体信息失败：%w", err)
	}

	if organizeByType {
		dstDir = dstDir.Join(GenMediaTypeFloderName(videoMeta.MediaType))
	}
	if organizeByCategory {
		dstDir = dstDir.Join(GenCategoryFloderName(info))
	}

	logrus.Infof("开始转移媒体文件：%s，目标目录：%s，转移类型：%s，组织方式：%t，分类方式：%t，刮削：%t",
		srcFile.String(), dstDir.String(), transferType, organizeByType, organizeByCategory, scrape)

	item, err := schemas.NewMediaItem(videoMeta, info)
	if err != nil {
		return fmt.Errorf("创建媒体项失败：%w", err)
	}

	var dst storage.StoragePath
	if scrape {
		dst, err = ArchiveMedia(srcFile, dstDir, transferType, item, info)
	} else {
		dst, err = ArchiveMedia(srcFile, dstDir, transferType, item, nil)
	}

	if err != nil {
		return fmt.Errorf("转移媒体文件失败：%w", err)
	}
	logrus.Infof("媒体文件转移成功：%s -> %s", srcFile.String(), dst.String())
	return nil
}
