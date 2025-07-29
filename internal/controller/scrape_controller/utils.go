package scrape_controller

import (
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
)

// 保存为 JPEG
func saveImageAsJPEG(img image.Image, imgFile *os.File) error {
	return jpeg.Encode(imgFile, img, &jpeg.Options{Quality: 90})
}

// 保存为 PNG
func saveImageAsPNG(img image.Image, imgFile *os.File) error {
	return png.Encode(imgFile, img)
}

// 保存图片到指定路径
// imgPath: 图片保存路径
// img: 要保存的图片
func SaveImage(imgPath string, img image.Image) error {
	imgFile, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	switch filepath.Ext(imgPath) {
	case ".jpg", ".jpeg":
		return saveImageAsJPEG(img, imgFile)
	default:
		return saveImageAsPNG(img, imgFile)
	}
}

// 判断文件或目录是否存在
func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 下载 TMDB 图片并保存到指定路径
// 自动根据 TMDB 图片的扩展名决定保存格式
// p: TMDB 中图片地址
// target: 目标路径，不带后缀名
func DownloadTMDBImageAndSave(p string, target string) error {
	target += path.Ext(p)
	if Exist(target) {
		return nil // 如果文件已存在，则不下载
	}

	img, err := tmdb_controller.DownloadImage(p)
	if err != nil {
		return err
	}
	return SaveImage(target, img)
}

// 下载 Fanart 图片并保存到指定路径
// 自动根据 Fanart 图片的扩展名决定保存格式
// url: Fanart 中图片地址
// target: 目标路径，不带后缀名
func DownloadFanartImageAndSave(url string, target string) error {
	target += path.Ext(url)
	if Exist(target) {
		return nil // 如果文件已存在，则不下载
	}

	img, err := fanart_controller.DownloadImage(url)
	if err != nil {
		return err
	}
	return SaveImage(target, img)
}
