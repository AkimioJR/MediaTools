package fanart

import (
	"fmt"
	"net/http"
)

type MovieImagesData struct {
	Name        string `json:"name"`
	TmdbID      string `json:"tmdb_id"`
	ImdbID      string `json:"imdb_id"`
	Hdmovielogo []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"hdmovielogo"`
	Moviedisc []struct {
		ID       string `json:"id"`
		URL      string `json:"url"`
		Lang     string `json:"lang"`
		Likes    string `json:"likes"`
		Disc     string `json:"disc"`
		DiscType string `json:"disc_type"`
	} `json:"moviedisc"`
	Movielogo []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"movielogo"`
	Movieposter []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"movieposter"`
	Hdmovieclearart []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"hdmovieclearart"`
	Movieart []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"movieart"`
	Moviebackground []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"moviebackground"`
	Moviebanner []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"moviebanner"`
	Moviethumb []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"moviethumb"`
}

// 获取电影的图片数据
// https://fanarttv.docs.apiary.io/#reference/movies/get-movies/get-images-for-movie
func (client *FanartClient) GetMovieImagesData(imdbID string) (*MovieImagesData, error) {
	var resp MovieImagesData
	err := client.DoRequest(
		http.MethodGet,
		"/movie/"+imdbID+"/images",
		nil,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewFanartError(err, fmt.Sprintf("获取电影「%s」图片数据失败：%v", imdbID, err))
	}
	return &resp, nil
}
