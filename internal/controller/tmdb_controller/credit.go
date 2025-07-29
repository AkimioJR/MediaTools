package tmdb_controller

import "MediaTools/internal/pkg/themoviedb/v3"

func GetMovieCredit(movieID int, language *string) (*themoviedb.MovieCredit, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetMovieCredit(movieID, language)
}

func GetTVSerieCredit(seriesID int, language *string) (*themoviedb.TVSerieCredit, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetTVSerieCredit(seriesID, language)
}

func GetTVSeasonCredit(seriesID int, seasonNumber int, language *string) (*themoviedb.TVSeasonCredit, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetTVSeasonCredit(seriesID, seasonNumber, language)
}

func GetTVEpisodeCredit(seriesID int, seasonNumber int, episodeNumber int, language *string) (*themoviedb.TVEpisodeCredit, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetTVEpisodeCredit(seriesID, seasonNumber, episodeNumber, language)
}
