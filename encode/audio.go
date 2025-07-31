package encode

import (
	"encoding/json"
	"strings"
)

// AudioEncode 音频编码类型枚举
type AudioEncode uint8

const (
	AudioEncodeUnknown AudioEncode = iota // 未知编码
	AudioEncodeAAC                        // AAC
	AudioEncodeAC3                        // AC-3
	AudioEncodeEAC3                       // E-AC-3/DD+
	AudioEncodeDTS                        // DTS
	AudioEncodeDTSHD                      // DTS-HD
	AudioEncodeDTSHDMA                    // DTS-HD MA
	AudioEncodeTrueHD                     // TrueHD
	AudioEncodeAtmos                      // Dolby Atmos
	AudioEncodeFLAC                       // FLAC
	AudioEncodeLPCM                       // LPCM
	AudioEncodeOpus                       // Opus
	AudioEncodeVorbis                     // Vorbis
	AudioEncodeMP3                        // MP3
)

// String 返回音频编码的字符串表示
func (ae AudioEncode) String() string {
	switch ae {
	case AudioEncodeAAC:
		return "AAC"
	case AudioEncodeAC3:
		return "AC-3"
	case AudioEncodeEAC3:
		return "E-AC-3"
	case AudioEncodeDTS:
		return "DTS"
	case AudioEncodeDTSHD:
		return "DTS-HD"
	case AudioEncodeDTSHDMA:
		return "DTS-HD MA"
	case AudioEncodeTrueHD:
		return "TrueHD"
	case AudioEncodeAtmos:
		return "Dolby Atmos"
	case AudioEncodeFLAC:
		return "FLAC"
	case AudioEncodeLPCM:
		return "LPCM"
	case AudioEncodeOpus:
		return "Opus"
	case AudioEncodeVorbis:
		return "Vorbis"
	case AudioEncodeMP3:
		return "MP3"
	default:
		return ""
	}
}

// ParseAudioEncode 从字符串解析音频编码
func ParseAudioEncode(s string) AudioEncode {
	s = strings.ToUpper(s)

	// 特殊处理DTS相关编码的复杂模式
	if strings.HasPrefix(s, "DTS") {
		// 标准化DTS格式
		normalized := s

		// 处理各种DTS-HD MA格式变体
		if strings.Contains(s, "MA") {
			// 匹配各种DTS-HD MA格式
			if strings.Contains(s, "HD") && strings.Contains(s, "MA") {
				return AudioEncodeDTSHDMA
			}
		}

		// 处理DTS-HD格式
		if strings.Contains(s, "HD") {
			return AudioEncodeDTSHD
		}

		// 纯DTS
		if normalized == "DTS" {
			return AudioEncodeDTS
		}
	}

	switch s {
	case "AAC", "AAC2.0", "AAC5.1":
		return AudioEncodeAAC
	case "AC3", "AC-3", "DD5.1":
		return AudioEncodeAC3
	case "EAC3", "E-AC-3", "DD+", "DD+7.1", "DDP", "DDP5.1", "DDP5", "DDP2.0":
		return AudioEncodeEAC3
	case "DTS":
		return AudioEncodeDTS
	case "DTSHD", "DTS-HD":
		return AudioEncodeDTSHD
	case "DTSHDMA", "DTS-HD MA", "DTS-HD MA5.1", "DTSMA":
		return AudioEncodeDTSHDMA
	case "TRUEHD", "TRUE-HD":
		return AudioEncodeTrueHD
	case "ATMOS":
		return AudioEncodeAtmos
	case "FLAC":
		return AudioEncodeFLAC
	case "LPCM":
		return AudioEncodeLPCM
	case "OPUS":
		return AudioEncodeOpus
	case "VORBIS":
		return AudioEncodeVorbis
	case "MP3":
		return AudioEncodeMP3
	default:
		return AudioEncodeUnknown
	}
}

func (ae AudioEncode) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ae.String() + `"`), nil
}

func (ae *AudioEncode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*ae = ParseAudioEncode(s)
	return nil
}
