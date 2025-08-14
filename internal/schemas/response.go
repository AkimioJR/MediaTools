package schemas

type ErrResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// 识别视频媒体信息的响应结构体
