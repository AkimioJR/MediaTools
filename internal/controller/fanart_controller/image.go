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

	data.Hdmovielogo = fanart.SortByLanguages(data.Hdmovielogo, config.Fanart.Languages)
	data.Movielogo = fanart.SortByLanguages(data.Movielogo, config.Fanart.Languages)
	data.Movieposter = fanart.SortByLanguages(data.Movieposter, config.Fanart.Languages)
	data.Hdmovieclearart = fanart.SortByLanguages(data.Hdmovieclearart, config.Fanart.Languages)
	data.Movieart = fanart.SortByLanguages(data.Movieart, config.Fanart.Languages)
	data.Moviebackground = fanart.SortByLanguages(data.Moviebackground, config.Fanart.Languages)
	data.Moviebanner = fanart.SortByLanguages(data.Moviebanner, config.Fanart.Languages)
	data.Moviethumb = fanart.SortByLanguages(data.Moviethumb, config.Fanart.Languages)

	return data, nil
}

func GetTVImagesData(thetvdbID int) (*fanart.TVImagesData, error) {
	lock.RLock()
	defer lock.RUnlock()

	data, err := client.GetTVImagesData(thetvdbID)
	if err != nil {
		return nil, err
	}

	data.Clearlogo = fanart.SortByLanguages(data.Clearlogo, config.Fanart.Languages)
	data.Hdtvlogo = fanart.SortByLanguages(data.Hdtvlogo, config.Fanart.Languages)
	data.Clearart = fanart.SortByLanguages(data.Clearart, config.Fanart.Languages)
	data.Showbackground = fanart.SortByLanguages(data.Showbackground, config.Fanart.Languages)
	data.Tvthumb = fanart.SortByLanguages(data.Tvthumb, config.Fanart.Languages)
	data.Seasonposter = fanart.SortByLanguages(data.Seasonposter, config.Fanart.Languages)
	data.Seasonthumb = fanart.SortByLanguages(data.Seasonthumb, config.Fanart.Languages)
	data.Hdclearart = fanart.SortByLanguages(data.Hdclearart, config.Fanart.Languages)
	data.Tvbanner = fanart.SortByLanguages(data.Tvbanner, config.Fanart.Languages)
	data.Characterart = fanart.SortByLanguages(data.Characterart, config.Fanart.Languages)
	data.Tvposter = fanart.SortByLanguages(data.Tvposter, config.Fanart.Languages)
	data.Seasonbanner = fanart.SortByLanguages(data.Seasonbanner, config.Fanart.Languages)

	return data, nil

}

func DownloadImage(url string) (image.Image, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.DownloadImage(url)
}
