package tmdb_controller

import "MediaTools/internal/pkg/themoviedb/v3"

func GetMovieCredits(movieID int, language *string) (*themoviedb.MovieCreditsResponse, error) {
	return client.GetMovieCredits(movieID, language)
}

func GetTVSeriesCredits(seriesID int, language *string) (*themoviedb.TVCreditsResponse, error) {
	return client.GetTVSeriesCredits(seriesID, language)
}

func GetTVSeasonCredits(seriesID int, seasonNumber int, language *string) (*themoviedb.TVCreditsResponse, error) {
	return client.GetTVSeasonCredits(seriesID, seasonNumber, language)
}

func GetTVEpisodeCredits(seriesID int, seasonNumber int, episodeNumber int, language *string) (*themoviedb.TVCreditsResponse, error) {
	return client.GetTVEpisodeCredits(seriesID, seasonNumber, episodeNumber, language)
}
