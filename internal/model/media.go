package model

import (
	"MediaTools/encode"
	"MediaTools/meta"
)

type MediaInfo struct {
	Title         string                 `json:"title"`          // 标题
	OriginalTitle string                 `json:"original_title"` // 原始标题
	Year          int                    `json:"year"`           // 年份
	MediaType     meta.MediaType         `json:"media_type"`     // 电影、电视剧
	TMDBID        uint64                 `json:"tmdb_id"`        // TMDB ID
	Part          string                 `json:"part"`           // 分段
	Version       uint8                  `json:"version"`        // 版本号
	ReleaseGroups []string               `json:"release_groups"` // 发布组
	Platform      meta.StreamingPlatform `json:"platform"`       // 流媒体平台

	// 资源相关信息
	ResourceType   meta.ResourceType     `json:"resource_type"`   // 资源类型
	ResourceEffect []meta.ResourceEffect `json:"resource_effect"` // 资源效果
	ResourcePix    meta.ResourcePix      `json:"resource_pix"`    // 分辨率
	VideoEncode    encode.VideoEncode    `json:"video_encode"`    // 视频编码
	AudioEncode    encode.AudioEncode    `json:"audio_encode"`    // 音频编码

	// 电视剧数据
	Season  string `json:"season"`  // 季
	Episode string `json:"episode"` // 集
}
