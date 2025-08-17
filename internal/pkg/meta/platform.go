package meta

import (
	"encoding/json"
	"strings"
)

// StreamingPlatforms 表示流媒体平台的唯一标识符（uint32的别名）
type StreamingPlatform uint32

// 平台ID枚举（示例，实际需要为每个平台分配唯一ID）
const (
	UnknownStreamingPlatform StreamingPlatform = iota
	Amazon
	Netflix
	AppleTV
	ITunes
	Disney
	Baha
	BiliBili
	BiliGlobal
	Crunchyroll
	YouTube
)

// String 返回流媒体平台的字符串表示
func (sp StreamingPlatform) String() string {
	switch sp {
	case Amazon:
		return "Amazon"
	case Netflix:
		return "Netflix"
	case AppleTV:
		return "Apple TV+"
	case ITunes:
		return "iTunes"
	case Disney:
		return "Disney+"
	case Baha:
		return "Baha"
	case BiliBili:
		return "BiliBili"
	case BiliGlobal:
		return "B-Global"
	case Crunchyroll:
		return "Crunchyroll"
	case YouTube:
		return "YouTube"
	default:
		return ""
	}
}

// ParseResourceType 从字符串解析资源类型
func (sp *StreamingPlatform) ParseString(s string) {
	switch strings.ToUpper(s) {
	case "AMAZON", "AMZN":
		*sp = Amazon
	case "NETFLIX", "NF":
		*sp = Netflix
	case "APPLE TV+", "ATVP":
		*sp = AppleTV
	case "ITUNES", "iT":
		*sp = ITunes
	case "DISNEY+", "DSNP":
		*sp = Disney
	case "BAHA":
		*sp = Baha
	case "BILIBILI", "BILI":
		*sp = BiliBili
	case "B-GLOBAL", "BG":
		*sp = BiliGlobal
	case "CRUNCHYROLL", "CR":
		*sp = Crunchyroll
	case "YOUTUBE", "YT":
		*sp = YouTube
	default:
		*sp = UnknownStreamingPlatform
	}
}

func (sp StreamingPlatform) MarshalJSON() ([]byte, error) {
	return []byte(`"` + sp.String() + `"`), nil
}

func (sp *StreamingPlatform) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	sp.ParseString(s)
	return nil
}
