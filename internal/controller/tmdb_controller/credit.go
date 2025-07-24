package tmdb_controller

import "MediaTools/internal/pkg/themoviedb/v3"

func GetMovieCredits(movieID int, language *string) (*themoviedb.MovieCredit, error) {
	return client.GetMovieCredit(movieID, language)
}

func GetTVSeriesCredits(seriesID int, language *string) (*themoviedb.TVSerieCredit, error) {
	return client.GetTVSerieCredit(seriesID, language)
}

func GetTVSeasonCredits(seriesID int, seasonNumber int, language *string) (*themoviedb.TVSeasonCredit, error) {
	return client.GetTVSeasonCredit(seriesID, seasonNumber, language)
}

func GetTVEpisodeCredits(seriesID int, seasonNumber int, episodeNumber int, language *string) (*themoviedb.TVEpisodeCredit, error) {
	return client.GetTVEpisodeCredit(seriesID, seasonNumber, episodeNumber, language)
}
