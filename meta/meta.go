package meta

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// MetaVideo 视频媒体信息结构体
type MetaVideo struct {
	// 基础信息
	title         string    // 标题
	orginalString string    // 原始字符串
	isFile        bool      // 是否是文件
	cnName        string    // 中文名
	enName        string    // 英文名
	year          uint      // 年份
	mediaType     MediaType // 媒体类型

	// 资源信息
	resourceType   ResourceType                // 来源/介质
	resourceEffect map[ResourceEffect]struct{} // 资源效果
	resourcePix    ResourcePix                 // 分辨率
	videoEncode    VideoEncode                 // 视频编码
	audioEncode    AudioEncode                 // 音频编码
	platform       StreamingPlatform           // 流媒体平台
	releaseGroups  []string                    // 发布组
	part           string                      // 分段
	// customization  string                      // 自定义词

	// 电视剧相关·
	beginSeason  *int // 起始季
	endSeason    *int // 结束集
	totalSeason  int  // 总季数
	beginEpisode *int // 起始集
	endEpisode   *int // 结束集
	totalEpisode int  // 总集数
}

func (meta *MetaVideo) GetTitle() string                               { return meta.title }          // GetSubtitle 获取标题
func (meta *MetaVideo) GetCNName() string                              { return meta.cnName }         // GetENName 获取中文名
func (meta *MetaVideo) GetENName() string                              { return meta.enName }         // GetENName 获取英文名
func (meta *MetaVideo) GetYear() uint                                  { return meta.year }           // GetYear 获取年份
func (meta *MetaVideo) GetType() MediaType                             { return meta.mediaType }      // MediaType
func (meta *MetaVideo) GetResourceType() ResourceType                  { return meta.resourceType }   // GetResourceType 获取资源类型
func (meta *MetaVideo) GetResourceEffect() map[ResourceEffect]struct{} { return meta.resourceEffect } // GetResourceEffect 获取资源效果
func (meta *MetaVideo) GetResourcePix() ResourcePix                    { return meta.resourcePix }    // GetResourcePix 获取资源分辨率
func (meta *MetaVideo) GetVideoEncode() VideoEncode                    { return meta.videoEncode }    // GetVideoEncode 获取视频编码
func (meta *MetaVideo) GetAudioEncode() AudioEncode                    { return meta.audioEncode }    // GetAudioEncode 获取音频编码
func (meta *MetaVideo) GetStreamingPlatform() StreamingPlatform        { return meta.platform }       // GetWebSource 获取网络来源
func (meta *MetaVideo) GetReleaseGroups() []string                     { return meta.releaseGroups }  // GetResourceTeam 获取资源组
func (meta *MetaVideo) GetPart() string                                { return meta.part }           // GetPart 获取分集信息

// 返回媒体的名字，有中文名优先返回中文名，否则返回英文名
func (meta *MetaVideo) GetName() string {
	cnName := meta.GetCNName()
	if cnName != "" {
		return cnName
	}
	return meta.GetENName()
}

func (meta *MetaVideo) GetResourceEffectStrings() []string {
	var effects []string
	for effect := range meta.resourceEffect {
		effects = append(effects, effect.String())
	}
	return effects
}

func (meta *MetaVideo) GetSeasons() []int {
	if meta.mediaType != MediaTypeTV {
		return []int{}
	}

	var seasons []int
	if meta.beginSeason != nil {
		if meta.endSeason != nil {
			for i := *meta.beginSeason; i <= *meta.endSeason; i++ {
				seasons = append(seasons, i)
			}
		} else {
			seasons = append(seasons, *meta.beginSeason)
		}
	}
	return seasons
}

func (meta *MetaVideo) GetSeasonStr() string {
	if meta.beginSeason == nil {
		return ""
	} else {
		if meta.endSeason == nil {
			return fmt.Sprintf("S%02d", *meta.beginSeason)
		}
		return fmt.Sprintf("S%02d-S%02d", *meta.beginSeason, *meta.endSeason)
	}
}

// GetEpisodes 获取集数列表
func (meta *MetaVideo) GetEpisodes() []int {
	if meta.mediaType != MediaTypeTV {
		return []int{}
	}

	var episodes []int
	if meta.beginEpisode != nil {
		if meta.endEpisode != nil {
			for i := *meta.beginEpisode; i <= *meta.endEpisode; i++ {
				episodes = append(episodes, i)
			}
		} else {
			episodes = append(episodes, *meta.beginEpisode)
		}
	}
	return episodes
}

func (meta *MetaVideo) GetEpisodeStr() string {
	if meta.beginEpisode == nil {
		return ""
	} else {
		if meta.endEpisode == nil {
			return fmt.Sprintf("E%02d", *meta.beginEpisode)
		}
		return fmt.Sprintf("E%02d-E%02d", *meta.beginEpisode, *meta.endEpisode)
	}
}

func ParseMetaVideo(originalString string, isFile bool) *MetaVideo {
	meta := &MetaVideo{
		orginalString:  originalString,
		isFile:         isFile,
		mediaType:      MediaTypeUnknown,
		resourceType:   ResourceTypeUnknown,
		resourceEffect: make(map[ResourceEffect]struct{}),
		releaseGroups:  findReleaseGroups(originalString), // 解析发布组
		platform:       UnknownStreamingPlatform,
	}

	title := originalString

	title = nameNoBeginRe.ReplaceAllString(title, "") // 去掉名称中第1个[]的内容
	loc := nameNoBeginRe.FindStringIndex(title)
	if loc != nil {
		title = title[:loc[0]] + title[loc[1]:]
	}
	title = yearRangeRe.ReplaceAllString(title, "${1}${2}") // 把xxxx-xxxx年份换成前一个年份，常出现在季集上
	title = fileSizeRe.ReplaceAllString(title, "")          // 把大小去掉
	title = dateFmtRe.ReplaceAllString(title, "")           // 把年月日去掉

	state := &parseState{
		tokens:         NewTokens(title), // 拆分tokens
		lastType:       lastTokenTypeUnknown,
		unknownNameStr: "",
		continueFlag:   true,
		stopNameFlag:   false,
		stopCNNameFlag: false,
	}
	for !state.tokens.isEnd() {
		state.tokens.GetNext() // 指向下一个
		state.continueFlag = true

		if state.continueFlag { // Part
			meta.parsePart(state)
		}

		if state.continueFlag { // 标题
			meta.parseName(state)
		}

		if state.continueFlag { // 年份
			meta.parseYear(state)
		}

		if state.continueFlag { // 辨率率
			meta.parseResourcePix(state)
		}

		if state.continueFlag { // 季度
			meta.parseSeason(state)
		}

		if state.continueFlag { // 集数
			meta.parseEpisode(state)
		}

		if state.continueFlag { // 资源类型
			meta.parseResourceType(state)
		}

		if state.continueFlag { // 流媒体平台
			meta.parsePlatform(state)
		}

		if state.continueFlag { // 视频编码
			meta.parseVideoEncode(state)
		}

		if state.continueFlag { // 音频编码
			meta.parseAudioEncode(state)
		}

	}

	meta.postProcess() // 后处理逻辑
	return meta
}

// 识别 Part
func (meta *MetaVideo) parsePart(s *parseState) {
	if meta.GetName() == "" {
		return
	}
	if meta.GetYear() == 0 &&
		meta.beginSeason == nil &&
		meta.beginEpisode == nil &&
		meta.GetResourcePix() == ResourcePixUnknown &&
		meta.GetResourceType() == ResourceTypeUnknown {
		return
	}

	token := s.tokens.Current()

	if partRe.MatchString(token) {
		meta.part = token

		nextToken := s.tokens.Peek()
		utf8Str := []rune(nextToken)
		length := len(utf8Str)
		if nextToken != "" {
			if (isDigits(nextToken) && (length == 1 || (length == 2 && utf8Str[0] == '0'))) ||
				contain([]string{"A", "B", "C", "I", "II", "III"}, strings.ToUpper(nextToken)) {
				meta.part += nextToken
			}
		}

		s.lastType = lastTokenTypePart
		s.continueFlag = false
	}
}

// 识别 CNName、ENName
func (meta *MetaVideo) parseName(s *parseState) {
	token := s.tokens.Current()

	if s.unknownNameStr != "" { // 回收标题
		if meta.cnName == "" {
			if meta.enName == "" {
				meta.enName = s.unknownNameStr
			} else if s.unknownNameStr != strconv.Itoa(int(meta.year)) {
				meta.enName += " " + s.unknownNameStr
			}
			s.lastType = lastTokenTypeEnName
		}
		s.unknownNameStr = ""
	}

	if s.stopNameFlag {
		return
	}

	if contain(meta.releaseGroups, token) { // 如果当前token是发布组，直接跳过
		s.continueFlag = false
		return
	}

	// 遇到AKA停止解析名称
	if strings.ToUpper(token) == "AKA" {
		s.continueFlag = false
		s.stopNameFlag = true
		return
	}

	// 遇到季集关键词，暂停处理
	if contain([]string{"共", "第", "季", "集", "话", "話", "期"}, token) {
		s.lastType = lastTokenTypeNameSeWords
		return
	}

	if isChinese(token) { // 中文处理
		// 含有中文，直接做为标题（连着的数字或者英文会保留），且不再取用后面出现的中文
		s.lastType = lastTokenTypeCnName
		if meta.cnName == "" {
			meta.cnName = token
		} else if !s.stopCNNameFlag {
			// 含有电影关键词或者不含特殊字符的中文可以继续拼接
			if contain([]string{"剧场版", "劇場版", "电影版", "電影版"}, token) ||
				(!nameNoChineseRe.MatchString(token) && !contain([]string{"共", "第", "季", "集", "话", "話", "期"}, token)) {
				meta.cnName += " " + token
			}
			s.stopCNNameFlag = true
		}
	} else { // 数字或罗马数字处理
		if isDigits(token) || isRomanNumeral(token) {
			if s.lastType == lastTokenTypeNameSeWords { // 如果前一个token是季集关键词，跳过
				return
			}

			if meta.GetName() != "" {
				if strings.HasPrefix(token, "0") { // 名字后面以0开头的不要，极有可能是集
					return
				}

				// 检查是否为数字
				if isDigits(token) {
					tokenInt, err := strconv.Atoi(token)
					if err != nil {
						return
					}

					if !isRomanNumeral(token) && s.lastType == lastTokenTypeCnName && tokenInt < 1900 { // 中文名后面跟的数字不是年份的极有可能是集
						return
					}
				}

				if (isDigits(token) && len(token) < 4) || isRomanNumeral(token) {
					// 4位以下的数字或者罗马数字，拼装到已有标题中
					switch s.lastType {
					case lastTokenTypeCnName:
						meta.cnName += " " + token
					case lastTokenTypeEnName:
						meta.enName += " " + token
					}
					s.continueFlag = false
				} else if isDigits(token) && len(token) == 4 { // 4位数字，可能是年份，也可能是标题的一部分
					if s.unknownNameStr == "" {
						s.unknownNameStr = token
					}
				}
			} else {
				// 名字未出现前的第一个数字，记下来
				if s.unknownNameStr == "" {
					s.unknownNameStr = token
				}
			}
		} else if seasonRe.MatchString(token) { // 季的处理
			if meta.enName != "" && strings.HasSuffix(strings.ToUpper(meta.enName), "SEASON") { // 如果匹配到季，英文名结尾为Season，说明Season属于标题，不应在后续作为干扰词去除
				meta.enName += " "
			}
			s.stopNameFlag = true
			return
		} else if episodeRe.MatchString(token) ||
			sourceRe.MatchString(token) ||
			effectRe.MatchString(token) ||
			resourcePixRe.MatchString(token) {
			s.stopNameFlag = true // 集、来源、版本等不要
			return
		} else {
			if isMediaExtension("." + strings.ToLower(token)) { // 后缀名不要
				return
			}

			// 英文或者英文+数字，拼装起来
			if meta.enName != "" {
				meta.enName += " " + token
			} else {
				meta.enName = token
			}
			s.lastType = lastTokenTypeEnName
		}
	}
}

// 识别年份
func (meta *MetaVideo) parseYear(s *parseState) {
	if meta.GetName() == "" {
		return
	}

	token := s.tokens.Current()
	if len([]rune(token)) != 4 {
		return
	}

	num, err := strconv.Atoi(token)
	if err != nil {
		return
	}
	if num < 1900 || num > 2050 {
		return
	}

	if meta.GetYear() != 0 {
		if meta.GetENName() != "" {
			meta.enName = strings.TrimSpace(meta.GetENName()) + " " + strconv.Itoa(int(meta.GetYear()))
		} else if meta.GetCNName() != "" {
			meta.cnName += " " + strconv.Itoa(int(meta.GetYear()))
		}
	} else if meta.GetENName() != "" && strings.HasSuffix(strings.ToLower(meta.GetENName()), "season") { // 如果匹配到年，且英文名结尾为Season，说明Season属于标题，不应在后续作为干扰词去除
		meta.enName += " "
	}

	meta.year = uint(num)
	s.lastType = lastTokenTypeYear
	s.continueFlag = false
	s.stopNameFlag = true
}

// 识别分辨率
func (meta *MetaVideo) parseResourcePix(s *parseState) {
	if meta.GetName() == "" {
		return
	}

	token := s.tokens.Current()

	// 使用第一个正则表达式匹配分辨率
	matches := resourcePixRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		s.lastType = lastTokenTypePix
		s.continueFlag = false
		s.stopNameFlag = true

		var resourcePixStr string
		// 遍历匹配的分组，找到非空的分组
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				resourcePixStr = matches[i]
				break
			}
		}

		if resourcePixStr != "" && meta.resourcePix == ResourcePixUnknown {
			// 如果分辨率是纯数字且不以k、p、i结尾，添加p后缀
			if isDigits(resourcePixStr) {
				lastChar := resourcePixStr[len(resourcePixStr)-1]
				if lastChar != 'k' && lastChar != 'p' && lastChar != 'i' {
					resourcePixStr = resourcePixStr + "p"
				}
			}
			meta.resourcePix = ParseResourcePix(strings.ToLower(resourcePixStr))
		}
		return
	}

	// 使用第二个正则表达式匹配分辨率
	matches2 := resourcePix2Re.FindStringSubmatch(token)
	if len(matches2) > 1 {
		s.lastType = lastTokenTypePix
		s.continueFlag = false
		s.stopNameFlag = true

		if meta.resourcePix == ResourcePixUnknown {
			meta.resourcePix = ParseResourcePix(strings.ToLower(matches2[1]))
		}
	}
}

// 识别季
func (meta *MetaVideo) parseSeason(s *parseState) {
	token := s.tokens.Current()

	if sxxexxRe.MatchString(token) { // 跳过 SxxExx 格式，让集数解析器处理
		return
	}

	// 使用季识别正则匹配
	matches := seasonRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		s.lastType = lastTokenTypeSeason
		meta.mediaType = MediaTypeTV
		s.stopNameFlag = true
		s.continueFlag = false

		// 从正则匹配结果中提取季数
		var seasonNum int
		var err error
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" && isDigits(matches[i]) {
				seasonNum, err = strconv.Atoi(matches[i])
				if err != nil {
					return
				}
				break
			}
		}

		if seasonNum > 0 {
			if meta.beginSeason == nil {
				meta.beginSeason = &seasonNum
				meta.totalSeason = 1
			} else {
				if seasonNum > *meta.beginSeason {
					meta.endSeason = &seasonNum
					meta.totalSeason = (seasonNum - *meta.beginSeason) + 1
					// 如果是文件且总季数大于1，重置结束季
					if meta.isFile && meta.totalSeason > 1 {
						meta.endSeason = nil
						meta.totalSeason = 1
					}
				}
			}
		}
	} else if isDigits(token) {
		// 检查是否为数字token
		tokenInt, err := strconv.Atoi(token)
		if err != nil {
			return
		}

		// 如果前一个token是SEASON且当前季为空且数字长度小于3
		if s.lastType == lastTokenTypeSeason &&
			meta.beginSeason == nil &&
			len(token) < 3 {
			meta.beginSeason = &tokenInt
			meta.totalSeason = 1
			s.lastType = lastTokenTypeSeason
			s.stopNameFlag = true
			s.continueFlag = false
			meta.mediaType = MediaTypeTV
		}
	} else if strings.ToUpper(token) == "SEASON" && meta.beginSeason == nil {
		// 遇到SEASON关键词
		s.lastType = lastTokenTypeSeason
	} else if meta.mediaType == MediaTypeTV && meta.beginSeason == nil {
		// 如果已确定为电视剧类型但没有季数，默认为第1季
		defaultSeason := 1
		meta.beginSeason = &defaultSeason
	}
}

// 识别集数
func (meta *MetaVideo) parseEpisode(s *parseState) {
	token := s.tokens.Current()

	sxxexxMatches := sxxexxRe.FindStringSubmatch(token) // 特殊处理 SxxExx 格式
	if len(sxxexxMatches) > 0 {                         // 同时解析季度和集数
		seasonStr := sxxexxMatches[1]
		episodeStr := sxxexxMatches[2]

		if seasonNum, err := strconv.Atoi(seasonStr); err == nil && seasonNum > 0 {
			if meta.beginSeason == nil {
				meta.beginSeason = &seasonNum
				meta.totalSeason = 1
			}
		}

		if episodeNum, err := strconv.Atoi(episodeStr); err == nil && episodeNum > 0 {
			if meta.beginEpisode == nil {
				meta.beginEpisode = &episodeNum
				meta.totalEpisode = 1
			}
		}

		s.lastType = lastTokenTypeEpisode
		s.continueFlag = false
		s.stopNameFlag = true
		meta.mediaType = MediaTypeTV
		return
	}

	// 使用集识别正则匹配
	matches := episodeRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		s.lastType = lastTokenTypeEpisode
		s.continueFlag = false
		s.stopNameFlag = true
		meta.mediaType = MediaTypeTV

		// 从正则匹配结果中提取集数
		var episodeNum int
		var err error
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" && isDigits(matches[i]) {
				episodeNum, err = strconv.Atoi(matches[i])
				if err != nil {
					return
				}
				break
			}
		}

		if episodeNum > 0 {
			if meta.beginEpisode == nil {
				meta.beginEpisode = &episodeNum
				meta.totalEpisode = 1
			} else {
				if episodeNum > *meta.beginEpisode {
					meta.endEpisode = &episodeNum
					meta.totalEpisode = (episodeNum - *meta.beginEpisode) + 1
					// 如果是文件且总集数大于2，重置结束集
					if meta.isFile && meta.totalEpisode > 2 {
						meta.endEpisode = nil
						meta.totalEpisode = 1
					}
				}
			}
		}
	} else if isDigits(token) {
		// 检查是否为数字token
		tokenInt, err := strconv.Atoi(token)
		if err != nil {
			return
		}
		length := len([]rune(token))

		if meta.beginEpisode != nil && // 情况1：已有起始集，没有结束集，且前一个token是episode
			meta.endEpisode == nil &&
			length < 5 &&
			tokenInt > *meta.beginEpisode &&
			s.lastType == lastTokenTypeEpisode {
			meta.endEpisode = &tokenInt
			meta.totalEpisode = (tokenInt - *meta.beginEpisode) + 1
			if meta.isFile && meta.totalEpisode > 2 {
				meta.endEpisode = nil
				meta.totalEpisode = 1
			}
			s.continueFlag = false
			meta.mediaType = MediaTypeTV
		} else if meta.beginEpisode == nil && // 情况2：没有起始集，数字长度在1-4之间，且不是年份或其他特殊token
			length > 1 && length < 4 &&
			s.lastType != lastTokenTypeYear &&
			s.lastType != lastTokenTypePix &&
			s.lastType != lastTokenTypeVideoEncode && // 避免将视频编码中的数字识别为集数
			token != s.unknownNameStr {

			meta.beginEpisode = &tokenInt
			meta.totalEpisode = 1
			s.lastType = lastTokenTypeEpisode
			s.continueFlag = false
			s.stopNameFlag = true
			meta.mediaType = MediaTypeTV
		} else if s.lastType == lastTokenTypeEpisode && // 情况3：前一个token是EPISODE关键词
			meta.beginEpisode == nil &&
			length < 5 {

			meta.beginEpisode = &tokenInt
			meta.totalEpisode = 1
			s.lastType = lastTokenTypeEpisode
			s.continueFlag = false
			s.stopNameFlag = true
			meta.mediaType = MediaTypeTV
		}
	} else if strings.ToUpper(token) == "EPISODE" { // 遇到EPISODE关键词
		s.lastType = lastTokenTypeEpisode
	}
}

// 识别资源类型
func (meta *MetaVideo) parseResourceType(s *parseState) {
	if meta.GetName() == "" {
		return
	}

	token := s.tokens.Current()
	tokenUpper := strings.ToUpper(token)

	// 处理特殊组合情况
	if tokenUpper == "DL" && s.lastType == lastTokenTypeSource && meta.resourceType == ResourceTypeWeb {
		meta.resourceType = ResourceTypeWebDL
		s.continueFlag = false
		return
	} else if tokenUpper == "RAY" && s.lastType == lastTokenTypeSource && meta.resourceType == ResourceTypeBlu {
		// UHD BluRay组合
		if meta.resourceType == ResourceTypeUHD {
			meta.resourceType = ResourceTypeUHDBluRay
		} else {
			meta.resourceType = ResourceTypeBluRay
		}
		s.continueFlag = false
		return
	} else if tokenUpper == "WEBDL" {
		meta.resourceType = ResourceTypeWebDL
		s.continueFlag = false
		return
	}

	// UHD REMUX组合
	if tokenUpper == "REMUX" && meta.resourceType == ResourceTypeBluRay {
		meta.resourceType = ResourceTypeBluRayRemux
		s.continueFlag = false
		return
	} else if tokenUpper == "BLURAY" && meta.resourceType == ResourceTypeUHD {
		meta.resourceType = ResourceTypeUHDBluRay
		s.continueFlag = false
		return
	}

	// 使用资源类型正则匹配
	matches := sourceRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		s.lastType = lastTokenTypeSource
		s.continueFlag = false
		s.stopNameFlag = true
		if meta.resourceType == ResourceTypeUnknown {
			meta.resourceType = ParseResourceType(matches[0])
		}
		return
	}
	// 使用效果正则匹配
	matches2 := effectRe.FindStringSubmatch(token)
	if len(matches2) > 0 {
		effect := ParseResourceEffect(matches2[0])
		if effect != ResourceEffectUnknown {
			s.lastType = lastTokenTypeEffect
			s.continueFlag = false
			s.stopNameFlag = true
			meta.resourceEffect[effect] = struct{}{}
		}
		return
	}
}

// 识别流媒体平台
func (meta *MetaVideo) parsePlatform(s *parseState) {
	// 检查是否已有名称
	if meta.GetName() == "" {
		return
	}

	token := s.tokens.Current()
	var platformName string
	queryRange := 1

	// 获取前一个token
	var prevToken string
	currentIndex := s.tokens.GetCurrentIndex()
	if currentIndex >= 2 {
		prevToken = s.tokens.GetByIndex(currentIndex - 2)
	}

	// 获取下一个token
	nextToken := s.tokens.Peek()

	// 检查当前token是否为流媒体平台
	p := ParseStreamingPlatform(token)
	if p != UnknownStreamingPlatform {
		platformName = p.String()
	} else {
		// 检查相邻token的组合
		adjacentTokens := []struct {
			token  string
			isNext bool
		}{
			{prevToken, false},
			{nextToken, true},
		}

		for _, adjacent := range adjacentTokens {
			if adjacent.token == "" || platformName != "" {
				continue
			}

			// 尝试不同的分隔符组合
			separators := []string{" ", "-"}
			for _, separator := range separators {
				var combinedToken string
				if adjacent.isNext {
					combinedToken = token + separator + adjacent.token
				} else {
					combinedToken = adjacent.token + separator + token
				}
				p := ParseStreamingPlatform(combinedToken)
				if p != UnknownStreamingPlatform {
					platformName = p.String()
					queryRange = 2
					if adjacent.isNext {
						s.tokens.GetNext() // 消费下一个token
					}
					break
				}
			}
			if platformName != "" {
				break
			}
		}
	}

	if platformName == "" {
		return
	}

	// 检查附近是否有WEB相关的token
	webTokens := []string{"WEB", "DL", "WEBDL", "WEBRIP"}
	matchStartIdx := currentIndex - queryRange
	matchEndIdx := currentIndex - 1
	startIndex := max(0, matchStartIdx-queryRange)
	endIndex := min(s.tokens.GetLength(), matchEndIdx+1+queryRange)

	tokensToCheck := s.tokens.GetTokensInRange(startIndex, endIndex)

	// 检查是否有WEB相关token
	hasWebToken := false
	for _, tok := range tokensToCheck {
		if tok != "" {
			upperTok := strings.ToUpper(tok)
			for _, webToken := range webTokens {
				if upperTok == webToken {
					hasWebToken = true
					break
				}
			}
			if hasWebToken {
				break
			}
		}
	}

	if hasWebToken {
		meta.platform = ParseStreamingPlatform(platformName)
		s.lastType = lastTokenTypePlatform
		s.continueFlag = false
		s.stopNameFlag = true
	}
}

// 识别视频编码
func (meta *MetaVideo) parseVideoEncode(s *parseState) {
	// 检查是否已有名称
	if meta.GetName() == "" {
		return
	}

	// 检查是否有其他必要信息
	if meta.GetYear() == 0 &&
		meta.resourcePix == ResourcePixUnknown &&
		meta.resourceType == ResourceTypeUnknown &&
		meta.beginSeason == nil &&
		meta.beginEpisode == nil {
		return
	}

	token := s.tokens.Current()
	tokenUpper := strings.ToUpper(token)

	// 使用视频编码正则表达式匹配
	matches := videoEncodeRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		s.continueFlag = false
		s.stopNameFlag = true
		s.lastType = lastTokenTypeVideoEncode

		if meta.videoEncode == VideoEncodeUnknown {
			// 从正则匹配结果中提取编码信息
			var encodeStr string
			for i := 1; i < len(matches); i++ {
				if matches[i] != "" {
					encodeStr = matches[i]
					break
				}
			}
			if encodeStr == "" {
				encodeStr = matches[0]
			}

			meta.videoEncode = ParseVideoEncode(encodeStr)
		}
		return
	}

	// 处理单独的 H 或 X 字母
	if tokenUpper == "H" || tokenUpper == "X" {
		s.continueFlag = false
		s.stopNameFlag = true
		s.lastType = lastTokenTypeVideoEncode
		return
	}

	// 处理 H264、H265、X264、X265 组合
	if (token == "264" || token == "265") &&
		s.lastType == lastTokenTypeVideoEncode {
		prevToken := ""
		currentIndex := s.tokens.GetCurrentIndex()
		if currentIndex >= 2 {
			prevToken = s.tokens.GetByIndex(currentIndex - 2)
		}
		prevTokenUpper := strings.ToUpper(prevToken)

		if prevTokenUpper == "H" || prevTokenUpper == "X" {
			encodeStr := prevTokenUpper + token
			meta.videoEncode = ParseVideoEncode(encodeStr)
			s.continueFlag = false
			s.stopNameFlag = true
		}
		return
	}

	// 处理 VC1、MPEG2 等数字组合
	if isDigits(token) && s.lastType == lastTokenTypeVideoEncode {
		prevToken := ""
		currentIndex := s.tokens.GetCurrentIndex()
		if currentIndex >= 2 {
			prevToken = s.tokens.GetByIndex(currentIndex - 2)
		}
		prevTokenUpper := strings.ToUpper(prevToken)

		if prevTokenUpper == "VC" || prevTokenUpper == "MPEG" {
			encodeStr := prevTokenUpper + token
			meta.videoEncode = ParseVideoEncode(encodeStr)
			s.continueFlag = false
			s.stopNameFlag = true
		}
		return
	}

	// 处理 10bit 编码
	if tokenUpper == "10BIT" {
		s.lastType = lastTokenTypeVideoEncode
		if meta.videoEncode == VideoEncodeUnknown { // 如果没有其他编码信息，设置为纯10bit编码
			meta.videoEncode = VideoEncode10bit
		} else { // 使用辅助方法升级为对应的10bit版本
			meta.videoEncode = meta.videoEncode.CombineWith10bit()
		}
		s.continueFlag = false
		s.stopNameFlag = true
		return
	}
}

// 识别音频编码
func (meta *MetaVideo) parseAudioEncode(s *parseState) {
	// 检查是否已有名称
	if meta.GetName() == "" {
		return
	}

	// 检查是否有其他必要信息
	if meta.GetYear() == 0 &&
		meta.resourcePix == ResourcePixUnknown &&
		meta.resourceType == ResourceTypeUnknown &&
		meta.beginSeason == nil &&
		meta.beginEpisode == nil {
		return
	}

	token := s.tokens.Current()

	// 使用音频编码正则表达式匹配
	matches := audioEncodeRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		s.continueFlag = false
		s.stopNameFlag = true
		s.lastType = lastTokenTypeAudioEncode

		matchedStr := matches[0]
		if meta.audioEncode == AudioEncodeUnknown {
			meta.audioEncode = ParseAudioEncode(matchedStr)
		} else {
			// 如果已有音频编码，进行组合处理
			currentEncodeStr := meta.audioEncode.String()
			newEncodeStr := ParseAudioEncode(matchedStr).String()

			if currentEncodeStr != "" && newEncodeStr != "" {
				// DTS相关编码使用 "-" 连接
				if strings.ToUpper(currentEncodeStr) == "DTS" {
					combinedStr := currentEncodeStr + "-" + newEncodeStr
					meta.audioEncode = ParseAudioEncode(combinedStr)
					if meta.audioEncode == AudioEncodeUnknown {
						// 如果无法解析组合编码，保持原有编码
						meta.audioEncode = ParseAudioEncode(currentEncodeStr)
					}
				} else {
					// 其他编码使用空格连接
					combinedStr := currentEncodeStr + " " + newEncodeStr
					meta.audioEncode = ParseAudioEncode(combinedStr)
					if meta.audioEncode == AudioEncodeUnknown {
						// 如果无法解析组合编码，保持原有编码
						meta.audioEncode = ParseAudioEncode(currentEncodeStr)
					}
				}
			}
		}
		return
	}

	// 处理数字token（用于音频编码的版本号等）
	if isDigits(token) && s.lastType == lastTokenTypeAudioEncode {
		if meta.audioEncode != AudioEncodeUnknown {
			currentStr := meta.audioEncode.String()
			if currentStr != "" {
				// 获取前一个token作为参考
				prevToken := ""
				currentIndex := s.tokens.GetCurrentIndex()
				if currentIndex >= 2 {
					prevToken = s.tokens.GetByIndex(currentIndex - 2)
				}

				var newEncodeStr string
				if isDigits(prevToken) {
					// 如果前一个token也是数字，用点号连接（如 7.1）
					newEncodeStr = currentStr + "." + token
				} else if len(currentStr) > 0 && currentStr[len(currentStr)-1] >= '0' && currentStr[len(currentStr)-1] <= '9' {
					// 如果当前编码末尾是数字，格式化为类似 "DTS 7.1" 的形式
					baseStr := currentStr[:len(currentStr)-1]
					lastDigit := string(currentStr[len(currentStr)-1])
					newEncodeStr = baseStr + " " + lastDigit + "." + token
				} else {
					// 普通情况，用空格连接
					newEncodeStr = currentStr + " " + token
				}

				// 尝试解析新的组合编码
				newEncode := ParseAudioEncode(newEncodeStr)
				if newEncode != AudioEncodeUnknown {
					meta.audioEncode = newEncode
				}
				// 注：如果解析失败，保持原有编码不变
			}
		}
		s.continueFlag = false
		return
	}
}

// postProcess 后处理：清理和优化解析结果
func (meta *MetaVideo) postProcess() {
	// 处理part
	if meta.part != "" && strings.ToUpper(meta.part) == "PART" {
		meta.part = ""
	}

	// 清理名称中的干扰字符
	meta.cnName = meta.fixName(meta.cnName)
	meta.enName = meta.fixName(meta.enName)

	// 英文名首字母大写
	if meta.enName != "" {
		words := strings.Fields(meta.enName)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			}
		}
		meta.enName = strings.Join(words, " ")
	}

	// 处理BluRay DIY标记
	if meta.resourceType != ResourceTypeUnknown &&
		(meta.resourceType == ResourceTypeBluRay ||
			meta.resourceType == ResourceTypeUHDBluRay ||
			meta.resourceType == ResourceTypeBluRayRemux) {
		// 检查原始字符串中是否包含DIY标记
		upperOriginal := strings.ToUpper(meta.orginalString)
		if strings.Contains(upperOriginal, "DIY") ||
			strings.Contains(upperOriginal, "-DIY@") {
			// 可以添加DIY标记到资源效果中
			meta.resourceEffect[ResourceEffectUnknown] = struct{}{} // 需要定义DIY效果类型
		}
	}
}

// fixName 清理名称中的干扰字符
func (meta *MetaVideo) fixName(name string) string {
	if name == "" {
		return name
	}

	// 使用干扰词过滤正则去除不需要的内容
	name = noStringRe.ReplaceAllString(name, "")
	name = strings.TrimSpace(name)

	// 合并多个空格为单个空格
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")

	// 如果名称是纯数字且小于1800，可能是误识别的集数
	if isDigits(name) {
		num, err := strconv.Atoi(name)
		if err == nil && num < 1800 &&
			meta.year == 0 &&
			meta.beginSeason == nil &&
			meta.resourcePix == ResourcePixUnknown &&
			meta.resourceType == ResourceTypeUnknown &&
			meta.audioEncode == AudioEncodeUnknown &&
			meta.videoEncode == VideoEncodeUnknown {

			// 如果还没有起始集，将此数字设为起始集
			if meta.beginEpisode == nil {
				meta.beginEpisode = &num
				meta.totalEpisode = 1
				meta.mediaType = MediaTypeTV
				return ""
			} else if meta.isInEpisode(num) && meta.beginSeason == nil {
				// 如果数字在集数范围内且没有季数信息，清空名称
				return ""
			}
		}
	}

	return name
}

// isInEpisode 检查数字是否在当前集数范围内
func (meta *MetaVideo) isInEpisode(episode int) bool {
	if meta.beginEpisode == nil {
		return false
	}

	if meta.endEpisode == nil {
		return episode == *meta.beginEpisode
	}

	return episode >= *meta.beginEpisode && episode <= *meta.endEpisode
}
