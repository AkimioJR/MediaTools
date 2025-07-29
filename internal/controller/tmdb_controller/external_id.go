package tmdb_controller

import "MediaTools/internal/pkg/themoviedb/v3"

func GetMovieExternalID(movieID int) (*themoviedb.MovieExternalID, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetMovieExternalID(movieID)
}

func GetTVSerieExternalID(tvID int) (*themoviedb.TVSerieExternalID, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetTVSerieExternalID(tvID)
}

func GetTVSeasonExternalID(tvID, seasonNumber int) (*themoviedb.TVSeasonExternalID, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetTVSeasonExternalID(tvID, seasonNumber)
}

func GetTVEpisodeExternalID(tvID, seasonNumber, episodeNumber int) (*themoviedb.TVEpisodeExternalID, error) {
	lock.RLock()
	defer lock.RUnlock()
	return client.GetTVEpisodeExternalID(tvID, seasonNumber, episodeNumber)
}
