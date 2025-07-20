package meta

import "strings"

// VideoEncode 视频编码类型枚举
type VideoEncode uint8

const (
	VideoEncodeUnknown    VideoEncode = iota // 未知编码
	VideoEncodeH264                          // H.264/AVC
	VideoEncodeH265                          // H.265/HEVC
	VideoEncodeH264_10bit                    // H.264 10bit
	VideoEncodeH265_10bit                    // H.265 10bit
	VideoEncodeAV1                           // AV1
	VideoEncodeAV1_10bit                     // AV1 10bit
	VideoEncodeXvid                          // Xvid
	VideoEncodeDivX                          // DivX
	VideoEncodeMPEG2                         // MPEG-2
	VideoEncodeMPEG4                         // MPEG-4
	VideoEncodeVC1                           // VC-1
	VideoEncodeAVS2                          // AVS2
	VideoEncodeAVS3                          // AVS3
	VideoEncode10bit                         // 纯10bit编码（未指定具体编码格式）
)

// String 返回视频编码的字符串表示
func (ve VideoEncode) String() string {
	switch ve {
	case VideoEncodeH264:
		return "H.264"
	case VideoEncodeH265:
		return "H.265"
	case VideoEncodeH264_10bit:
		return "H.264 10bit"
	case VideoEncodeH265_10bit:
		return "H.265 10bit"
	case VideoEncodeAV1:
		return "AV1"
	case VideoEncodeAV1_10bit:
		return "AV1 10bit"
	case VideoEncodeXvid:
		return "Xvid"
	case VideoEncodeDivX:
		return "DivX"
	case VideoEncodeMPEG2:
		return "MPEG-2"
	case VideoEncodeMPEG4:
		return "MPEG-4"
	case VideoEncodeVC1:
		return "VC-1"
	case VideoEncodeAVS2:
		return "AVS2"
	case VideoEncodeAVS3:
		return "AVS3"
	case VideoEncode10bit:
		return "10bit"
	default:
		return ""
	}
}

// ParseVideoEncode 从字符串解析视频编码
func ParseVideoEncode(s string) VideoEncode {
	switch strings.ToUpper(s) {
	case "H264", "H.264", "X264", "AVC":
		return VideoEncodeH264
	case "H265", "H.265", "X265", "HEVC":
		return VideoEncodeH265
	case "H264 10BIT", "H.264 10BIT", "X264 10BIT", "AVC 10BIT":
		return VideoEncodeH264_10bit
	case "H265 10BIT", "H.265 10BIT", "X265 10BIT", "HEVC 10BIT":
		return VideoEncodeH265_10bit
	case "AV1":
		return VideoEncodeAV1
	case "AV1 10BIT":
		return VideoEncodeAV1_10bit
	case "10BIT":
		return VideoEncode10bit
	case "XVID":
		return VideoEncodeXvid
	case "DIVX":
		return VideoEncodeDivX
	case "MPEG2", "MPEG-2":
		return VideoEncodeMPEG2
	case "MPEG4", "MPEG-4":
		return VideoEncodeMPEG4
	case "VC1", "VC-1":
		return VideoEncodeVC1
	case "AVS2":
		return VideoEncodeAVS2
	case "AVS3":
		return VideoEncodeAVS3
	default:
		return VideoEncodeUnknown
	}
}

// CombineWith10bit 将现有编码与10bit组合，返回对应的10bit版本
func (ve VideoEncode) CombineWith10bit() VideoEncode {
	switch ve {
	case VideoEncodeH264:
		return VideoEncodeH264_10bit
	case VideoEncodeH265:
		return VideoEncodeH265_10bit
	case VideoEncodeAV1:
		return VideoEncodeAV1_10bit
	default:
		// 如果是未知编码或其他不支持10bit的编码，返回纯10bit
		return VideoEncode10bit
	}
}

// Is10bit 检查编码是否为10bit版本
func (ve VideoEncode) Is10bit() bool {
	switch ve {
	case VideoEncodeH264_10bit, VideoEncodeH265_10bit, VideoEncodeAV1_10bit, VideoEncode10bit:
		return true
	default:
		return false
	}
}

// GetBaseEncode 获取编码的基础版本（去除10bit标识）
func (ve VideoEncode) GetBaseEncode() VideoEncode {
	switch ve {
	case VideoEncodeH264_10bit:
		return VideoEncodeH264
	case VideoEncodeH265_10bit:
		return VideoEncodeH265
	case VideoEncodeAV1_10bit:
		return VideoEncodeAV1
	case VideoEncode10bit:
		return VideoEncodeUnknown
	default:
		return ve
	}
}

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
	switch strings.ToUpper(s) {
	case "AAC":
		return AudioEncodeAAC
	case "AC3", "AC-3":
		return AudioEncodeAC3
	case "EAC3", "E-AC-3", "DD+", "DDP":
		return AudioEncodeEAC3
	case "DTS":
		return AudioEncodeDTS
	case "DTSHD", "DTS-HD":
		return AudioEncodeDTSHD
	case "DTSHDMA", "DTS-HD MA", "DTSMA":
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
