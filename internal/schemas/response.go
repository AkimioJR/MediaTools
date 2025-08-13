package schemas

type ErrResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// 识别视频媒体信息的响应结构体
type RecognizeMediaResponse struct {
	Item       *MediaItem `json:"item"`        // 识别到的媒体项
	CustomRule string     `json:"custom_rule"` //应用自定义规则
	MetaRule   string     `json:"meta_rule"`   // 应用的媒体规则
}
