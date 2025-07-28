package tmdb_controller

import (
	"MediaTools/internal/pkg/themoviedb/v3"
	"image"
)

func GetImageURL(path string) string {
	return client.GetImageURL(path)
}

func DownloadImage(path string) (image.Image, error) {
	return client.DownloadImage(path)
}

func GetMovieImage(movieID int) (*themoviedb.MovieImage, error) {
	return client.GetMovieImage(movieID, nil, nil)
}

func GetTVSerieImage(tvID int) (*themoviedb.TVSerieImage, error) {
	return client.GetTVSerieImage(tvID, nil, nil)
}

func GetTVSeasonImage(tvID, seasonNumber int) (*themoviedb.TVSeasonImage, error) {
	return client.GetTVSeasonImage(tvID, seasonNumber, nil)
}

func GetTVEpisodeImage(tvID, seasonNumber, episodeNumber int) (*themoviedb.TVEpisodeImage, error) {
	return client.GetTVEpisodeImage(tvID, seasonNumber, episodeNumber, nil)
}
