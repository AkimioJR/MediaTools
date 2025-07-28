package fanart_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/fanart/v3"
	"image"
)

func GetMovieImagesData(imdbID string) (*fanart.MovieImagesData, error) {
	lock.RLock()
	defer lock.RUnlock()

	data, err := client.GetMovieImagesData(imdbID)
	if err != nil {
		return nil, err
	}

	data.HDMovieLogo = fanart.SortByLanguages(data.HDMovieLogo, config.Fanart.Languages)
	data.MovieLogo = fanart.SortByLanguages(data.MovieLogo, config.Fanart.Languages)
	data.MoviePoster = fanart.SortByLanguages(data.MoviePoster, config.Fanart.Languages)
	data.HDMovieClearArt = fanart.SortByLanguages(data.HDMovieClearArt, config.Fanart.Languages)
	data.MovieArt = fanart.SortByLanguages(data.MovieArt, config.Fanart.Languages)
	data.MovieBackground = fanart.SortByLanguages(data.MovieBackground, config.Fanart.Languages)
	data.MovieBanner = fanart.SortByLanguages(data.MovieBanner, config.Fanart.Languages)
	data.MovieThumb = fanart.SortByLanguages(data.MovieThumb, config.Fanart.Languages)

	return data, nil
}

func GetTVImagesData(thetvdbID int) (*fanart.TVImagesData, error) {
	lock.RLock()
	defer lock.RUnlock()

	data, err := client.GetTVImagesData(thetvdbID)
	if err != nil {
		return nil, err
	}

	data.ClearLogo = fanart.SortByLanguages(data.ClearLogo, config.Fanart.Languages)
	data.HDTVLogo = fanart.SortByLanguages(data.HDTVLogo, config.Fanart.Languages)
	data.ClearArt = fanart.SortByLanguages(data.ClearArt, config.Fanart.Languages)
	data.ShowBackground = fanart.SortByLanguages(data.ShowBackground, config.Fanart.Languages)
	data.TVThumb = fanart.SortByLanguages(data.TVThumb, config.Fanart.Languages)
	data.SeasonPoster = fanart.SortByLanguages(data.SeasonPoster, config.Fanart.Languages)
	data.SeasonThumb = fanart.SortByLanguages(data.SeasonThumb, config.Fanart.Languages)
	data.HDClearArt = fanart.SortByLanguages(data.HDClearArt, config.Fanart.Languages)
	data.TVBanner = fanart.SortByLanguages(data.TVBanner, config.Fanart.Languages)
	data.CharacterArt = fanart.SortByLanguages(data.CharacterArt, config.Fanart.Languages)
	data.TVPoster = fanart.SortByLanguages(data.TVPoster, config.Fanart.Languages)
	data.SeasonBanner = fanart.SortByLanguages(data.SeasonBanner, config.Fanart.Languages)

	return data, nil

}

func DownloadImage(url string) (image.Image, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.DownloadImage(url)
}
