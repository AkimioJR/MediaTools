package fanart

import (
	"fmt"
	"net/http"
	"net/url"
)

type MovieImagesData struct {
	Name            string     `json:"name"`
	TmdbID          string     `json:"tmdb_id"`
	ImdbID          string     `json:"imdb_id"`
	HDMovieLogo     []BaseInfo `json:"hdmovielogo"`
	MovieDisc       []DiscInfo `json:"moviedisc"`
	MovieLogo       []BaseInfo `json:"movielogo"`
	MoviePoster     []BaseInfo `json:"movieposter"`
	HDMovieClearArt []BaseInfo `json:"hdmovieclearart"`
	MovieArt        []BaseInfo `json:"movieart"`
	MovieBackground []BaseInfo `json:"moviebackground"`
	MovieBanner     []BaseInfo `json:"moviebanner"`
	MovieThumb      []BaseInfo `json:"moviethumb"`
}

// 获取电影的图片数据
// https://fanarttv.docs.apiary.io/#reference/movies/get-movies/get-images-for-movie
func (client *FanartClient) GetMovieImagesData(imdbID string) (*MovieImagesData, error) {
	var resp MovieImagesData
	err := client.DoRequest(
		http.MethodGet,
		"/movies/"+imdbID,
		url.Values{},
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewFanartError(fmt.Sprintf("获取电影「%s」图片数据失败", imdbID), err)
	}
	return &resp, nil
}
