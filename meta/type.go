package meta

import (
	"regexp"
)

// MediaType 媒体类型枚举
type MediaType uint8

const (
	MediaTypeUnknown MediaType = iota // 未知
	MediaTypeMovie                    //电影
	MediaTypeTV                       //电视剧
)

var (
	// 基础预处理正则
	nameNoBeginRe = regexp.MustCompile(`^[\[【].+?[\]】]`)
	yearRangeRe   = regexp.MustCompile(`([\s.]+)(\d{4})-(\d{4})`)
	fileSizeRe    = regexp.MustCompile(`(?i)[0-9.]+\s*[MGT]i?B\b`)
	dateFmtRe     = regexp.MustCompile(`\d{4}[\s._-]\d{1,2}[\s._-]\d{1,2}`)

	// 季集识别正则
	seasonRe = regexp.MustCompile(`(?i)S(\d{3})|^S(\d{1,3})$|S(\d{1,3})E|S(\d{1,3})$`)
	// seasonFullRe = regexp.MustCompile(`^Season\s+(\d{1,3})$|^S(\d{1,3})$`)
	episodeRe = regexp.MustCompile(`(?i)EP?(\d{2,4})$|^EP?(\d{1,4})$|^S\d{1,2}EP?(\d{1,4})$|S\d{2}EP?(\d{2,4})|S\d{2}E(\d{2,4})`)
	sxxexxRe  = regexp.MustCompile(`(?i)^S(\d{1,3})E(\d{1,4})$`)

	// Part识别正则
	partRe = regexp.MustCompile(`(?i)(^PART[0-9ABI]{0,2}$|^CD[0-9]{0,2}$|^DVD[0-9]{0,2}$|^DISK[0-9]{0,2}$|^DISC[0-9]{0,2}$)`)

	// 资源类型识别正则
	sourceRe = regexp.MustCompile(`(?i)^BLURAY$|^HDTV$|^UHDTV$|^HDDVD$|^WEBRIP$|^DVDRIP$|^BDRIP$|^BLU$|^WEB$|^BD$|^HDRip$|^REMUX$|^UHD$`)
	effectRe = regexp.MustCompile(`(?i)^SDR$|^HDR\d*$|^DOLBY$|^DOVI$|^DV$|^3D$|^REPACK$|^HLG$|^HDR10(\+|Plus)$|^EDR$|^HQ$`)

	// 分辨率识别正则
	resourcePixBaseRe     = regexp.MustCompile(`(?i)^[SBUHD]*((480|576|720|1080|1440|2160|2880|4320)[PI]*)|^[SBUHD]*(\d{3,4}[PI]+)|\d{3,4}X(\d{3,4})`)
	resourcePixStandardRe = regexp.MustCompile(`(?i)(^[248]+[KPI])`)

	// 编码识别正则
	videoEncodeRe = regexp.MustCompile(`(?i)^(H26[45])$|^(x26[45])$|^AVC$|^HEVC$|^VC\d?$|^MPEG\d?$|^Xvid$|^DivX$|^AV1$|^HDR\d*$|^AVS(\+|[23])$`)
	audioEncodeRe = regexp.MustCompile(`(?i)^DTS\d?$|^DTSHD$|^DTSHDMA$|^Atmos$|^TrueHD\d?$|^AC3$|^\dAudios?$|^DDP\d?$|^DD\+\d?$|^DD\d?$|^LPCM\d?$|^AAC\d?$|^FLAC\d?$|^HD\d?$|^MA\d?$|^HR\d?$|^Opus\d?$|^Vorbis\d?$|^AV[3S]A$`)

	// 干扰词过滤正则
	noStringRe = regexp.MustCompile(`(?i)^PTS|^AOD|^CHC|^[A-Z]{1,4}TV[\-0-9UVHDK]*` +
		`|HBO$|\s+HBO|\d{1,2}th|\d{1,2}bit|NETFLIX|AMAZON|IMAX|^3D|\s+3D|^BBC\s+|\s+BBC|BBC$|DISNEY\+?|XXX|\s+DC$` +
		`|[第\s共]+[0-9一二三四五六七八九十\-\s]+季` +
		`|[第\s共]+[0-9一二三四五六七八九十百零\-\s]+[集话話]` +
		`|连载|日剧|美剧|电视剧|动画片|动漫|欧美|西德|日韩|超高清|高清|无水印|下载|蓝光|翡翠台|梦幻天堂·龙网|★?\d*月?新番` +
		`|最终季|合集|[多中国英葡法俄日韩德意西印泰台港粤双文语简繁体特效内封官译外挂]+字幕|版本|出品|台版|港版|\w+字幕组|\w+字幕社` +
		`|未删减版|UNCUT$|UNRATE$|WITH EXTRAS$|RERIP$|SUBBED$|PROPER$|REPACK$|SEASON$|EPISODE$|Complete$|Extended$|Extended Version$` +
		`|S\d{2}\s*-\s*S\d{2}|S\d{2}|\s+S\d{1,2}|EP?\d{2,4}\s*-\s*EP?\d{2,4}|EP?\d{2,4}|\s+EP?\d{1,4}` +
		`|CD[\s.]*[1-9]|DVD[\s.]*[1-9]|DISK[\s.]*[1-9]|DISC[\s.]*[1-9]` +
		`|[248]K|\d{3,4}[PIX]+` +
		`|CD[\s.]*[1-9]|DVD[\s.]*[1-9]|DISK[\s.]*[1-9]|DISC[\s.]*[1-9]|\s+GB`)

	// 中文名过滤正则
	nameNoChineseRe = regexp.MustCompile(`.*版|.*字幕`)
	// nameMovieWordsRe = regexp.MustCompile(`(?i)movie|film|cinema`)
	// nameSeWordsRe    = regexp.MustCompile(`(?i)season|series|episode|ep|se`)

	// Token分割正则
	tokenSplitRe = regexp.MustCompile(`\.|\s+|\(|\)|\[|]|-|【|】|/|～|;|&|\||#|_|「|」|~`)
)

type lastTokenType uint8

const (
	lastTokenTypeUnknown     lastTokenType = iota // 未知
	lastTokenTypecntitle                          // 中文名
	lastTokenTypeentitle                          // 英文名
	lastTokenTypeYear                             // 年份
	lastTokenTypeSeason                           // 季
	lastTokenTypeEpisode                          // 集
	lastTokenTypePart                             // 识别Part
	lastTokenTypeNameSeWords                      // 季集关键词
	lastTokenTypePix                              // 分辨率
	lastTokenTypeSource                           // 资源类型
	lastTokenTypeEffect                           // 效果类型
	lastTokenTypePlatform                         // 流媒体平台
	lastTokenTypeVideoEncode                      // 视频编码
	lastTokenTypeAudioEncode                      // 音频编码
)

type parseState struct {
	tokens          *Tokens
	lastType        lastTokenType
	unknownNameStr  string // 回收标题
	stopNameFlag    bool
	stopcntitleFlag bool
	continueFlag    bool
}
