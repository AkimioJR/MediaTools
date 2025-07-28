package fanart

import (
	"fmt"
	"net/http"
	"strconv"
)

type TVImagesData struct {
	Name           string        `json:"name"`
	ThetvdbID      string        `json:"thetvdb_id"`
	Clearlogo      []Image       `json:"clearlogo"`
	Hdtvlogo       []Image       `json:"hdtvlogo"`
	Clearart       []Image       `json:"clearart"`
	Showbackground []SeasonImage `json:"showbackground"`
	Tvthumb        []Image       `json:"tvthumb"`
	Seasonposter   []Image       `json:"seasonposter"`
	Seasonthumb    []SeasonImage `json:"seasonthumb"`
	Hdclearart     []Image       `json:"hdclearart"`
	Tvbanner       []Image       `json:"tvbanner"`
	Characterart   []Image       `json:"characterart"`
	Tvposter       []Image       `json:"tvposter"`
	Seasonbanner   []SeasonImage `json:"seasonbanner"`
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
