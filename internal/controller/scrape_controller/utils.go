package scrape_controller

import (
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/schemas/storage"
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

// 保存图片到指定路径
// imgPath: 图片保存路径
// img: 要保存的图片
func SaveImage(imgFile storage.StoragePath, img image.Image) error {
	var (
		buff bytes.Buffer
		err  error
	)
	switch filepath.Ext(imgFile.GetPath()) {
	case ".jpg", ".jpeg": // 保存为 JPEG
		err = png.Encode(&buff, img)
	default: // 保存为 PNG
		err = jpeg.Encode(&buff, img, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return err
	}
	return storage_controller.CreateFile(imgFile, &buff)
}

// 下载 TMDB 图片并保存到指定路径
// 自动根据 TMDB 图片的扩展名决定保存格式
// p: TMDB 中图片地址
// target: 目标路径，不带后缀名
func DownloadTMDBImageAndSave(p string, target string, storageType storage.StorageType) error {
	target += filepath.Ext(p)
	dstFile, err := storage_controller.GetPath(target, storageType)
	if err != nil {
		return err
	}
	exists, err := storage_controller.Exists(dstFile)
	if err != nil {
		return err
	}
	if exists {
		return nil // 如果文件已存在，则跳过下载
	}

	img, err := tmdb_controller.DownloadImage(p)
	if err != nil {
		return err
	}
	return SaveImage(dstFile, img)
}

// 下载 Fanart 图片并保存到指定路径
// 自动根据 Fanart 图片的扩展名决定保存格式
// url: Fanart 中图片地址
// target: 目标路径，不带后缀名
func DownloadFanartImageAndSave(url string, target string, storageType storage.StorageType) error {
	target += filepath.Ext(url)
	dstFile, err := storage_controller.GetPath(target, storageType)
	if err != nil {
		return err
	}
	exists, err := storage_controller.Exists(dstFile)
	if err != nil {
		return err
	}
	if exists {
		return nil // 如果文件已存在，则不下载
	}

	img, err := fanart_controller.DownloadImage(url)
	if err != nil {
		return err
	}
	return SaveImage(dstFile, img)
}

func bytes2Reader(p []byte) (io.Reader, error) {
	var buffer bytes.Buffer
	_, err := buffer.Write(p)
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}

var supportImgExts = []string{".jpg", ".jpeg", ".png"}

func getSupportImage(arr []string) string {
	for i := range arr {
		ext := strings.ToLower(filepath.Ext(arr[i]))
		if slices.Contains(supportImgExts, ext) {
			return arr[i]
		}
	}
	return ""
}
