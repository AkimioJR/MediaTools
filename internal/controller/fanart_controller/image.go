package fanart_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/fanart/v3"
	"encoding/json"
	"image"

	"github.com/sirupsen/logrus"
)

func GetMovieImagesData(imdbID string) (*fanart.MovieImagesData, error) {
	lock.RLock()
	defer lock.RUnlock()

	data, err := client.GetMovieImagesData(imdbID)
	if err != nil {
		return nil, err
	}

	fanart.SortByLanguages(data.HDMovieLogo, config.Fanart.Languages)
	fanart.SortByLanguages(data.MovieLogo, config.Fanart.Languages)
	fanart.SortByLanguages(data.MoviePoster, config.Fanart.Languages)
	fanart.SortByLanguages(data.HDMovieClearArt, config.Fanart.Languages)
	fanart.SortByLanguages(data.MovieArt, config.Fanart.Languages)
	fanart.SortByLanguages(data.MovieBackground, config.Fanart.Languages)
	fanart.SortByLanguages(data.MovieBanner, config.Fanart.Languages)
	fanart.SortByLanguages(data.MovieThumb, config.Fanart.Languages)

	return data, nil
}

func GetTVImagesData(thetvdbID int) (*fanart.TVImagesData, error) {
	lock.RLock()
	defer lock.RUnlock()

	data, err := client.GetTVImagesData(thetvdbID)
	if err != nil {
		return nil, err
	}
	defer func() {
		d, _ := json.Marshal(data)
		logrus.Debugf("获取 Fanart TV 数据: %s", string(d))
	}()

	fanart.SortByLanguages(data.ClearLogo, config.Fanart.Languages)
	fanart.SortByLanguages(data.HDTVLogo, config.Fanart.Languages)
	fanart.SortByLanguages(data.ClearArt, config.Fanart.Languages)
	fanart.SortByLanguages(data.ShowBackground, config.Fanart.Languages)
	fanart.SortByLanguages(data.TVThumb, config.Fanart.Languages)
	fanart.SortByLanguages(data.SeasonPoster, config.Fanart.Languages)
	fanart.SortByLanguages(data.SeasonThumb, config.Fanart.Languages)
	fanart.SortByLanguages(data.HDClearArt, config.Fanart.Languages)
	fanart.SortByLanguages(data.TVBanner, config.Fanart.Languages)
	fanart.SortByLanguages(data.CharacterArt, config.Fanart.Languages)
	fanart.SortByLanguages(data.TVPoster, config.Fanart.Languages)
	fanart.SortByLanguages(data.SeasonBanner, config.Fanart.Languages)

	return data, nil

}

func DownloadImage(url string) (image.Image, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.DownloadImage(url)
}
