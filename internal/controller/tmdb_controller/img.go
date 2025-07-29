package tmdb_controller

import (
	"MediaTools/internal/pkg/themoviedb/v3"
	"image"
)

func GetImageURL(path string) string {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetImageURL(path)
}

func DownloadImage(path string) (image.Image, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.DownloadImage(path)
}

func GetMovieImage(movieID int) (*themoviedb.MovieImage, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetMovieImage(movieID, nil, nil)
}

func GetTVSerieImage(tvID int) (*themoviedb.TVSerieImage, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetTVSerieImage(tvID, nil, nil)
}

func GetTVSeasonImage(tvID, seasonNumber int) (*themoviedb.TVSeasonImage, error) {
	lock.RLock()
	defer lock.RUnlock()
	return client.GetTVSeasonImage(tvID, seasonNumber, nil)
}

func GetTVEpisodeImage(tvID, seasonNumber, episodeNumber int) (*themoviedb.TVEpisodeImage, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetTVEpisodeImage(tvID, seasonNumber, episodeNumber, nil)
}
