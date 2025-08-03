package meta

import (
	"encoding/json"
	"strings"
)

// ResourcePix 分辨率类型枚举
type ResourcePix uint8

const (
	ResourcePixUnknown ResourcePix = iota // 未知分辨率
	ResourcePix480p                       // 480p
	ResourcePix720p                       // 720p
	ResourcePix1080p                      // 1080p
	ResourcePix1440p                      // 1440p (2K)
	ResourcePix2160p                      // 2160p (4K)
	ResourcePix4320p                      // 4320p (8K)
)

// String 返回分辨率的字符串表示
func (rp ResourcePix) String() string {
	switch rp {
	case ResourcePix480p:
		return "480p"
	case ResourcePix720p:
		return "720p"
	case ResourcePix1080p:
		return "1080p"
	case ResourcePix1440p:
		return "1440p"
	case ResourcePix2160p:
		return "2160p"
	case ResourcePix4320p:
		return "4320p"
	default:
		return ""
	}
}

// ParseResourcePix 从字符串解析分辨率
func ParseResourcePix(s string) ResourcePix {
	switch strings.ToLower(s) {
	case "480p", "480i":
		return ResourcePix480p
	case "720p", "720i", "hd":
		return ResourcePix720p
	case "1080p", "1080i", "fullhd":
		return ResourcePix1080p
	case "1440p", "2k":
		return ResourcePix1440p
	case "2160p", "4k", "uhd":
		return ResourcePix2160p
	case "4320p", "8k":
		return ResourcePix4320p
	default:
		return ResourcePixUnknown
	}
}

func (rp ResourcePix) MarshalJSON() ([]byte, error) {
	return []byte(`"` + rp.String() + `"`), nil
}

func (rp *ResourcePix) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*rp = ParseResourcePix(s)
	return nil
}

// ResourceEffect 资源效果枚举
type ResourceEffect uint8

const (
	ResourceEffectUnknown   ResourceEffect = iota // 未知
	ResourceEffectSDR                             // SDR
	ResourceEffectHDR                             // HDR
	ResourceEffectHDR10                           // HDR10
	ResourceEffectHDR10Plus                       // HDR10+
	ResourceEffectDolby                           // Dolby
	ResourceEffectDovi                            // DOVI
	ResourceEffectDV                              // DV
	ResourceEffect3D                              // 3D
	ResourceEffectRepack                          // REPACK
	ResourceEffectHLG                             // HLG
	ResourceEffectEDR                             // EDR
	ResourceEffectHQ                              // HQ
)

// String 返回资源效果的字符串表示
func (re ResourceEffect) String() string {
	switch re {
	case ResourceEffectSDR:
		return "SDR"
	case ResourceEffectHDR:
		return "HDR"
	case ResourceEffectHDR10:
		return "HDR10"
	case ResourceEffectHDR10Plus:
		return "HDR10+"
	case ResourceEffectDolby:
		return "Dolby"
	case ResourceEffectDovi:
		return "DOVI"
	case ResourceEffectDV:
		return "DV"
	case ResourceEffect3D:
		return "3D"
	case ResourceEffectRepack:
		return "REPACK"
	case ResourceEffectHLG:
		return "HLG"
	case ResourceEffectEDR:
		return "EDR"
	case ResourceEffectHQ:
		return "HQ"
	default:
		return ""
	}
}

// ParseResourceEffect 从字符串解析资源效果
func ParseResourceEffect(s string) ResourceEffect {
	switch strings.ToUpper(s) {
	case "SDR":
		return ResourceEffectSDR
	case "HDR":
		return ResourceEffectHDR
	case "HDR10":
		return ResourceEffectHDR10
	case "HDR10+", "HDR10PLUS":
		return ResourceEffectHDR10Plus
	case "DOLBY":
		return ResourceEffectDolby
	case "DOVI":
		return ResourceEffectDovi
	case "DV":
		return ResourceEffectDV
	case "3D":
		return ResourceEffect3D
	case "REPACK":
		return ResourceEffectRepack
	case "HLG":
		return ResourceEffectHLG
	case "EDR":
		return ResourceEffectEDR
	case "HQ":
		return ResourceEffectHQ
	default:
		return ResourceEffectUnknown
	}
}

func (re ResourceEffect) MarshalJSON() ([]byte, error) {
	return []byte(`"` + re.String() + `"`), nil
}

func (re *ResourceEffect) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*re = ParseResourceEffect(s)
	return nil
}

// ResourceType 资源类型枚举
type ResourceType uint8

const (
	ResourceTypeUnknown     ResourceType = iota // 未知
	ResourceTypeBluRay                          // BluRay
	ResourceTypeHDTV                            // HDTV
	ResourceTypeUHDTV                           // UHDTV
	ResourceTypeHDDVD                           // HDDVD
	ResourceTypeWebRip                          // WebRip
	ResourceTypeDVDRip                          // DVDRip
	ResourceTypeBDRip                           // BDRip
	ResourceTypeBlu                             // Blu
	ResourceTypeWeb                             // Web
	ResourceTypeBD                              // BD
	ResourceTypeHDRip                           // HDRip
	ResourceTypeRemux                           // Remux
	ResourceTypeUHD                             // UHD
	ResourceTypeWebDL                           // WEB-DL
	ResourceTypeUHDBluRay                       // UHD BluRay
	ResourceTypeBluRayRemux                     // BluRay REMUX
)

// String 返回资源类型的字符串表示
func (rt ResourceType) String() string {
	switch rt {
	case ResourceTypeBluRay:
		return "BluRay"
	case ResourceTypeHDTV:
		return "HDTV"
	case ResourceTypeUHDTV:
		return "UHDTV"
	case ResourceTypeHDDVD:
		return "HDDVD"
	case ResourceTypeWebRip:
		return "WebRip"
	case ResourceTypeDVDRip:
		return "DVDRip"
	case ResourceTypeBDRip:
		return "BDRip"
	case ResourceTypeBlu:
		return "Blu"
	case ResourceTypeWeb:
		return "Web"
	case ResourceTypeBD:
		return "BD"
	case ResourceTypeHDRip:
		return "HDRip"
	case ResourceTypeRemux:
		return "Remux"
	case ResourceTypeUHD:
		return "UHD"
	case ResourceTypeWebDL:
		return "WEB-DL"
	case ResourceTypeUHDBluRay:
		return "UHD BluRay"
	case ResourceTypeBluRayRemux:
		return "BluRay REMUX"
	default:
		return ""
	}
}

func (rt ResourceType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + rt.String() + `"`), nil
}

func (rt *ResourceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*rt = ParseResourceType(s)
	return nil
}

// ParseResourceType 从字符串解析资源类型
func ParseResourceType(s string) ResourceType {
	switch strings.ToUpper(s) {
	case "BLURAY":
		return ResourceTypeBluRay
	case "HDTV":
		return ResourceTypeHDTV
	case "UHDTV":
		return ResourceTypeUHDTV
	case "HDDVD":
		return ResourceTypeHDDVD
	case "WEBRIP":
		return ResourceTypeWebRip
	case "DVDRIP":
		return ResourceTypeDVDRip
	case "BDRIP":
		return ResourceTypeBDRip
	case "BLU":
		return ResourceTypeBlu
	case "WEB":
		return ResourceTypeWeb
	case "BD":
		return ResourceTypeBD
	case "HDRIP":
		return ResourceTypeHDRip
	case "REMUX":
		return ResourceTypeRemux
	case "UHD":
		return ResourceTypeUHD
	case "WEB-DL", "WEBDL":
		return ResourceTypeWebDL
	case "UHD BLURAY":
		return ResourceTypeUHDBluRay
	case "BLURAY REMUX":
		return ResourceTypeBluRayRemux
	default:
		return ResourceTypeUnknown
	}
}
