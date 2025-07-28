package scrape_controller

import (
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

// 下载图片并保存到指定路径
// 自动根据 TMDB 图片的扩展名决定保存格式
// p: TMDB 中图片地址
// target: 目标路径，不带后缀名
func DownloadImageAndSave(p string, target string) error {
	img, err := tmdb_controller.DownloadImage(p)
	if err != nil {
		return err
	}
	target += path.Ext(p)
	return SaveImage(target, img)
}
