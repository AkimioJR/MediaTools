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
	ClearLogo      []BaseInfo   `json:"clearlogo"`
	HDTVLogo       []BaseInfo   `json:"hdtvlogo"`
	ClearArt       []BaseInfo   `json:"clearart"`
	ShowBackground []SeasonInfo `json:"showbackground"`
	TVThumb        []BaseInfo   `json:"tvthumb"`
	SeasonPoster   []BaseInfo   `json:"seasonposter"`
	SeasonThumb    []SeasonInfo `json:"seasonthumb"`
	HDClearArt     []BaseInfo   `json:"hdclearart"`
	TVBanner       []BaseInfo   `json:"tvbanner"`
	CharacterArt   []BaseInfo   `json:"characterart"`
	TVPoster       []BaseInfo   `json:"tvposter"`
	SeasonBanner   []SeasonInfo `json:"seasonbanner"`
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
		return nil, NewFanartError(fmt.Sprintf("获取剧集「%d」图片数据失败", thetvdbID), err)
	}
	return &resp, nil
}
