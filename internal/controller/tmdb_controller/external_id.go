package tmdb_controller

import (
	"MediaTools/internal/pkg/themoviedb/v3"

	"github.com/sirupsen/logrus"
)

func GetMovieExternalID(movieID int) (*themoviedb.MovieExternalID, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电影（TMDB ID: %d）外部ID", movieID)
	return client.GetMovieExternalID(movieID)
}

func GetTVSerieExternalID(tvID int) (*themoviedb.TVSerieExternalID, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）外部ID", tvID)
	return client.GetTVSerieExternalID(tvID)
}

func GetTVSeasonExternalID(tvID, seasonNumber int) (*themoviedb.TVSeasonExternalID, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）S%02d外部ID", tvID, seasonNumber)
	return client.GetTVSeasonExternalID(tvID, seasonNumber)
}

func GetTVEpisodeExternalID(tvID, seasonNumber, episodeNumber int) (*themoviedb.TVEpisodeExternalID, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）S%02dE%02d外部ID", tvID, seasonNumber, episodeNumber)
	return client.GetTVEpisodeExternalID(tvID, seasonNumber, episodeNumber)
}
