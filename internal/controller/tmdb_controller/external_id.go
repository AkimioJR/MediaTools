package tmdb_controller

import "MediaTools/internal/pkg/themoviedb/v3"

func GetMovieExternalID(movieID int) (*themoviedb.MovieExternalID, error) {
	return client.GetMovieExternalID(movieID)
}

func GetTVSerieExternalID(tvID int) (*themoviedb.TVSerieExternalID, error) {
	return client.GetTVSerieExternalID(tvID)
}

func GetTVSeasonExternalID(tvID, seasonNumber int) (*themoviedb.TVSeasonExternalID, error) {
	return client.GetTVSeasonExternalID(tvID, seasonNumber)
}

func GetTVEpisodeExternalID(tvID, seasonNumber, episodeNumber int) (*themoviedb.TVEpisodeExternalID, error) {
	return client.GetTVEpisodeExternalID(tvID, seasonNumber, episodeNumber)
}
