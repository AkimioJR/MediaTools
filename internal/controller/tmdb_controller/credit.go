package tmdb_controller

import (
	"MediaTools/internal/pkg/themoviedb/v3"

	"github.com/sirupsen/logrus"
)

func GetMovieCredit(movieID int, language *string) (*themoviedb.MovieCredit, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电影（TMDB ID: %d）演职员信息", movieID)
	return client.GetMovieCredit(movieID, language)
}

func GetTVSerieCredit(seriesID int, language *string) (*themoviedb.TVSerieCredit, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）演职员信息", seriesID)
	return client.GetTVSerieCredit(seriesID, language)
}

func GetTVSeasonCredit(seriesID int, seasonNumber int, language *string) (*themoviedb.TVSeasonCredit, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）第 %d 季演职员信息", seriesID, seasonNumber)
	return client.GetTVSeasonCredit(seriesID, seasonNumber, language)
}

func GetTVEpisodeCredit(seriesID int, seasonNumber int, episodeNumber int, language *string) (*themoviedb.TVEpisodeCredit, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）第 %d 季第 %d 集演职员信息", seriesID, seasonNumber, episodeNumber)
	return client.GetTVEpisodeCredit(seriesID, seasonNumber, episodeNumber, language)
}
