package fanart

import (
	"fmt"
	"net/http"
	"strconv"
)

type TVImagesData struct {
	Name      string `json:"name"`
	ThetvdbID string `json:"thetvdb_id"`
	Clearlogo []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"clearlogo"`
	Hdtvlogo []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"hdtvlogo"`
	Clearart []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"clearart"`
	Showbackground []struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Lang   string `json:"lang"`
		Likes  string `json:"likes"`
		Season string `json:"season"`
	} `json:"showbackground"`
	Tvthumb []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"tvthumb"`
	Seasonposter []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"seasonposter"`
	Seasonthumb []struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Lang   string `json:"lang"`
		Likes  string `json:"likes"`
		Season string `json:"season"`
	} `json:"seasonthumb"`
	Hdclearart []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"hdclearart"`
	Tvbanner []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"tvbanner"`
	Characterart []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"characterart"`
	Tvposter []struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Lang  string `json:"lang"`
		Likes string `json:"likes"`
	} `json:"tvposter"`
	Seasonbanner []struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Lang   string `json:"lang"`
		Likes  string `json:"likes"`
		Season string `json:"season"`
	} `json:"seasonbanner"`
}

// 获取剧集的图片数据
// https://fanarttv.docs.apiary.io/#reference/tv/get-show/get-images-for-show
func (client *FanartClient) GetTVImagesData(thetvdbID int) (*TVImagesData, error) {
	var resp TVImagesData
	err := client.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(thetvdbID)+"/images",
		nil,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewFanartError(err, fmt.Sprintf("获取剧集「%d」图片数据失败：%v", thetvdbID, err))
	}
	return &resp, nil
}
