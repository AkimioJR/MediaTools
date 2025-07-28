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
	Hdmovielogo     []BaseInfo `json:"hdmovielogo"`
	Moviedisc       []DiscInfo `json:"moviedisc"`
	Movielogo       []BaseInfo `json:"movielogo"`
	Movieposter     []BaseInfo `json:"movieposter"`
	Hdmovieclearart []BaseInfo `json:"hdmovieclearart"`
	Movieart        []BaseInfo `json:"movieart"`
	Moviebackground []BaseInfo `json:"moviebackground"`
	Moviebanner     []BaseInfo `json:"moviebanner"`
	Moviethumb      []BaseInfo `json:"moviethumb"`
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
		return nil, NewFanartError(err, fmt.Sprintf("获取电影「%s」图片数据失败：%v", imdbID, err))
	}
	return &resp, nil
}
