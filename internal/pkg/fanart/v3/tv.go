package fanart

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVImagesData struct {
	Name           string       `json:"name"`
	ThetvdbID      string       `json:"thetvdb_id"`
	Clearlogo      []BaseInfo   `json:"clearlogo"`
	Hdtvlogo       []BaseInfo   `json:"hdtvlogo"`
	Clearart       []BaseInfo   `json:"clearart"`
	Showbackground []SeasonInfo `json:"showbackground"`
	Tvthumb        []BaseInfo   `json:"tvthumb"`
	Seasonposter   []BaseInfo   `json:"seasonposter"`
	Seasonthumb    []SeasonInfo `json:"seasonthumb"`
	Hdclearart     []BaseInfo   `json:"hdclearart"`
	Tvbanner       []BaseInfo   `json:"tvbanner"`
	Characterart   []BaseInfo   `json:"characterart"`
	Tvposter       []BaseInfo   `json:"tvposter"`
	Seasonbanner   []SeasonInfo `json:"seasonbanner"`
}

// 获取剧集的图片数据
// https://fanarttv.docs.apiary.io/#reference/tv/get-show/get-images-for-show
func (client *FanartClient) GetTVImagesData(thetvdbID int) (*TVImagesData, error) {
	var resp TVImagesData
	err := client.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(thetvdbID)+"/images",
		url.Values{},
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewFanartError(err, fmt.Sprintf("获取剧集「%d」图片数据失败：%v", thetvdbID, err))
	}
	return &resp, nil
}
