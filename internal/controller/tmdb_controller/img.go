package tmdb_controller

import "image"

func GetImageURL(path string) string {
	return client.GetImageURL(path)
}

func DownloadImage(path string) (image.Image, error) {
	return client.DownloadImage(path)
}
