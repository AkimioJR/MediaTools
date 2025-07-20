package meta

import (
	"MediaTools/encode"
	"MediaTools/utils"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// MetaVideo 视频媒体信息结构体
type MetaVideo struct {
	// 基础信息
	orginalTitle   string    // 原始标题
	processedTitle string    // 处理后的标题
	isFile         bool      // 是否是媒体文件
	cntitle        string    // 中文标题
	entitle        string    // 英文标题
	year           uint      // 年份
	mediaType      MediaType // 媒体类型

	// 资源信息
	resourceType   ResourceType                // 来源/介质
	resourceEffect map[ResourceEffect]struct{} // 资源效果
	resourcePix    ResourcePix                 // 分辨率
	videoEncode    encode.VideoEncode          // 视频编码
	audioEncode    encode.AudioEncode          // 音频编码
	platform       StreamingPlatform           // 流媒体平台
	releaseGroups  []string                    // 发布组
	part           string                      // 分段
	version        uint8                       // 版本号
	// customization  string                      // 自定义词

	// 电视剧相关·
	beginSeason  *int // 起始季
	endSeason    *int // 结束集
	totalSeason  int  // 总季数
	beginEpisode *int // 起始集
	endEpisode   *int // 结束集
	totalEpisode int  // 总集数
}

func (meta *MetaVideo) GetCNTitle() string                             { return meta.cntitle }        // GetCNTitle 获取中文标题
func (meta *MetaVideo) GetENTitle() string                             { return meta.entitle }        // GetENTitle 获取英文标题
func (meta *MetaVideo) GetYear() uint                                  { return meta.year }           // GetYear 获取年份
func (meta *MetaVideo) GetType() MediaType                             { return meta.mediaType }      // MediaType
func (meta *MetaVideo) GetResourceType() ResourceType                  { return meta.resourceType }   // GetResourceType 获取资源类型
func (meta *MetaVideo) GetResourceEffect() map[ResourceEffect]struct{} { return meta.resourceEffect } // GetResourceEffect 获取资源效果
func (meta *MetaVideo) GetResourcePix() ResourcePix                    { return meta.resourcePix }    // GetResourcePix 获取资源分辨率
func (meta *MetaVideo) GetVideoEncode() encode.VideoEncode             { return meta.videoEncode }    // GetVideoEncode 获取视频编码
func (meta *MetaVideo) GetAudioEncode() encode.AudioEncode             { return meta.audioEncode }    // GetAudioEncode 获取音频编码
func (meta *MetaVideo) GetStreamingPlatform() StreamingPlatform        { return meta.platform }       // GetWebSource 获取网络来源
func (meta *MetaVideo) GetReleaseGroups() []string                     { return meta.releaseGroups }  // GetResourceTeam 获取资源组
func (meta *MetaVideo) GetPart() string                                { return meta.part }           // GetPart 获取分集信息
func (meta *MetaVideo) GetVersion() uint8                              { return meta.version }        // GetVersion 获取版本号，若未识别到版本号，则返回1

// 获取标题
// 有中文标题优先返回中文标题
// 否则返回英文标题
func (meta *MetaVideo) GetTitle() string {
	if meta.cntitle != "" {
		return meta.cntitle
	}
	return meta.entitle
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

func ParseMetaVideo(title string) *MetaVideo {
	meta := &MetaVideo{
		orginalTitle:   title,
		mediaType:      MediaTypeUnknown,
		resourceType:   ResourceTypeUnknown,
		resourceEffect: make(map[ResourceEffect]struct{}),
		releaseGroups:  findReleaseGroups(title), // 解析发布组
		platform:       UnknownStreamingPlatform,
		version:        ParseVersion(title), // 解析版本号
	}

	if utils.IsMediaExtension(path.Ext(title)) {
		title = strings.TrimSuffix(title, path.Ext(title)) // 去掉文件扩展名
		meta.isFile = true
	}

	loc := nameNoBeginRe.FindStringIndex(title) // 去掉名称中第1个[]的内容（一般是发布组）
	if loc != nil {
		title = title[:loc[0]] + title[loc[1]:]
	}
	title = yearRangeRe.ReplaceAllString(title, "${1}${2}") // 把xxxx-xxxx年份换成前一个年份，常出现在季集上
	title = fileSizeRe.ReplaceAllString(title, "")          // 把大小去掉
	title = dateFmtRe.ReplaceAllString(title, "")           // 把年月日去掉
	title = strings.TrimSpace(title)                        // 去掉首尾空格
	meta.processedTitle = title

	state := &parseState{
		tokens:          NewTokens(title), // 拆分tokens
		lastType:        lastTokenTypeUnknown,
		unknownNameStr:  "",
		continueFlag:    true,
		stopNameFlag:    false,
		stopcntitleFlag: false,
	}
	for !state.tokens.IsEnd() {
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
	if meta.GetTitle() == "" {
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
			if (utils.IsDigits(nextToken) && (length == 1 || (length == 2 && utf8Str[0] == '0'))) ||
				slices.Contains([]string{"A", "B", "C", "I", "II", "III"}, strings.ToUpper(nextToken)) {
				meta.part += nextToken
			}
		}

		s.lastType = lastTokenTypePart
		s.continueFlag = false
	}
}

// 识别 cntitle、entitle
func (meta *MetaVideo) parseName(s *parseState) {
	token := s.tokens.Current()

	if s.unknownNameStr != "" { // 回收标题
		if meta.cntitle == "" {
			if meta.entitle == "" {
				meta.entitle = s.unknownNameStr
			} else if s.unknownNameStr != strconv.Itoa(int(meta.year)) {
				meta.entitle += " " + s.unknownNameStr
			}
			s.lastType = lastTokenTypeentitle
		}
		s.unknownNameStr = ""
	}

	if s.stopNameFlag {
		return
	}

	if slices.Contains(meta.releaseGroups, token) { // 如果当前token是发布组，直接跳过
		s.continueFlag = false
		return
	}

	// 遇到季集关键词，暂停处理
	if slices.Contains([]string{"共", "第", "季", "集", "话", "話", "期"}, token) {
		s.lastType = lastTokenTypeNameSeWords
		return
	}

	if utils.IsChinese(token) { // 中文处理
		// 含有中文，直接做为标题（连着的数字或者英文会保留），且不再取用后面出现的中文
		s.lastType = lastTokenTypecntitle
		if meta.cntitle == "" {
			meta.cntitle = token
		} else if !s.stopcntitleFlag {
			// 含有电影关键词或者不含特殊字符的中文可以继续拼接
			if slices.Contains([]string{"剧场版", "劇場版", "电影版", "電影版"}, token) ||
				(!nameNoChineseRe.MatchString(token) && !slices.Contains([]string{"共", "第", "季", "集", "话", "話", "期"}, token)) {
				meta.cntitle += " " + token
			}
			s.stopcntitleFlag = true
		}
	} else { // 数字或罗马数字处理
		if utils.IsDigits(token) || utils.IsRomanNumeral(token) {
			if s.lastType == lastTokenTypeNameSeWords { // 如果前一个token是季集关键词，跳过
				return
			}

			if meta.GetTitle() != "" {
				if strings.HasPrefix(token, "0") { // 名字后面以0开头的不要，极有可能是集
					return
				}

				// 检查是否为数字
				if utils.IsDigits(token) {
					tokenInt, err := strconv.Atoi(token)
					if err != nil {
						return
					}

					if !utils.IsRomanNumeral(token) && s.lastType == lastTokenTypecntitle && tokenInt < 1900 { // 中文名后面跟的数字不是年份的极有可能是集
						return
					}
				}

				if (utils.IsDigits(token) && len(token) < 4) || utils.IsRomanNumeral(token) {
					// 4位以下的数字或者罗马数字，需要判断是否应该拼装到标题中
					if utils.IsDigits(token) {
						if tokenInt, err := strconv.Atoi(token); err == nil {
							// 如果数字在标题的合理范围内(如系列编号)，则附加到标题
							// 但是排除典型的集数范围(1-99)，除非有特殊情况
							if len(token) <= 2 && tokenInt >= 1 && tokenInt <= 99 {
								// 检查标题内容，判断是否更可能是标题的一部分
								currentTitle := ""
								if s.lastType == lastTokenTypecntitle {
									currentTitle = meta.cntitle
								} else if s.lastType == lastTokenTypeentitle {
									currentTitle = meta.entitle
								}

								// 如果标题包含"Part", "Season", "第"等关键词，数字更可能是标题的一部分
								titleLower := strings.ToLower(currentTitle)
								if strings.Contains(titleLower, "part") ||
									strings.Contains(titleLower, "season") ||
									strings.Contains(currentTitle, "第") ||
									strings.Contains(currentTitle, "系列") {
									// 附加到标题
									switch s.lastType {
									case lastTokenTypecntitle:
										meta.cntitle += " " + token
									case lastTokenTypeentitle:
										meta.entitle += " " + token
									}
									s.continueFlag = false
									return
								}

								// 如果数字过大（如 >100），不太可能是简单的集数，更可能是标题的一部分
								if tokenInt > 100 {
									switch s.lastType {
									case lastTokenTypecntitle:
										meta.cntitle += " " + token
									case lastTokenTypeentitle:
										meta.entitle += " " + token
									}
									s.continueFlag = false
									return
								}

								// 否则可能是集数，不附加到标题，让后续处理
								return
							}
							// 3位数字通常是标题的一部分（如电影名称）
							if len(token) == 3 {
								switch s.lastType {
								case lastTokenTypecntitle:
									meta.cntitle += " " + token
								case lastTokenTypeentitle:
									meta.entitle += " " + token
								}
								s.continueFlag = false
								return
							}
						}
					}
					// 其他情况（罗马数字或其他特殊数字）拼装到已有标题中
					switch s.lastType {
					case lastTokenTypecntitle:
						meta.cntitle += " " + token
					case lastTokenTypeentitle:
						meta.entitle += " " + token
					}
					s.continueFlag = false
				} else if utils.IsDigits(token) && len(token) == 4 { // 4位数字，可能是年份，也可能是标题的一部分
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
			if meta.entitle != "" && strings.HasSuffix(strings.ToUpper(meta.entitle), "SEASON") { // 如果匹配到季，英文名结尾为Season，说明Season属于标题，不应在后续作为干扰词去除
				meta.entitle += " "
			}
			s.stopNameFlag = true
			return
		} else if episodeRe.MatchString(token) ||
			sourceRe.MatchString(token) ||
			effectRe.MatchString(token) ||
			resourcePixBaseRe.MatchString(token) {
			s.stopNameFlag = true // 集、来源、版本等不要
			return
		} else {
			if utils.IsMediaExtension("." + strings.ToLower(token)) { // 后缀名不要
				return
			}

			// 英文或者英文+数字，拼装起来
			if meta.entitle != "" {
				meta.entitle += " " + token
			} else {
				meta.entitle = token
			}
			s.lastType = lastTokenTypeentitle
		}
	}
}

// 识别年份
func (meta *MetaVideo) parseYear(s *parseState) {
	if meta.GetTitle() == "" {
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
		if meta.GetENTitle() != "" {
			meta.entitle = strings.TrimSpace(meta.GetENTitle()) + " " + strconv.Itoa(int(meta.GetYear()))
		} else if meta.GetCNTitle() != "" {
			meta.cntitle += " " + strconv.Itoa(int(meta.GetYear()))
		}
	} else if meta.GetENTitle() != "" && strings.HasSuffix(strings.ToLower(meta.GetENTitle()), "season") { // 如果匹配到年，且英文名结尾为Season，说明Season属于标题，不应在后续作为干扰词去除
		meta.entitle += " "
	}

	meta.year = uint(num)
	s.lastType = lastTokenTypeYear
	s.continueFlag = false
	s.stopNameFlag = true
}

// 识别分辨率
func (meta *MetaVideo) parseResourcePix(s *parseState) {
	if meta.GetTitle() == "" {
		return
	}

	token := s.tokens.Current()
	var r ResourcePix = ResourcePixUnknown

	// 先使用更精确的2K/4K/8K匹配
	matches := resourcePixStandardRe.FindStringSubmatch(token)
	if len(matches) > 1 {
		if meta.resourcePix == ResourcePixUnknown {
			r = ParseResourcePix(matches[1])
			if r != ResourcePixUnknown {
				goto match // 如果匹配到有效的分辨率
			}
		}
	}

	// 再使用通用分辨率匹配
	matches = resourcePixBaseRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				resourcePixStr := matches[i]
				if resourcePixStr != "" && meta.resourcePix == ResourcePixUnknown {
					if utils.IsDigits(resourcePixStr) { // 如果分辨率是纯数字且不以k、p、i结尾，添加p后缀
						lastChar := resourcePixStr[len(resourcePixStr)-1]
						if lastChar != 'k' && lastChar != 'p' && lastChar != 'i' {
							resourcePixStr = resourcePixStr + "p"
						}
					}
					r = ParseResourcePix(resourcePixStr)
					if r != ResourcePixUnknown {
						goto match // 如果匹配到有效的分辨率
					}
				}
			}
		}

	}

	return // 没有匹配到有效的分辨率

match: // 如果匹配到有效的分辨率
	s.lastType = lastTokenTypePix
	s.continueFlag = false
	s.stopNameFlag = true
	meta.resourcePix = r
}

// 识别季
func (meta *MetaVideo) parseSeason(s *parseState) {
	token := s.tokens.Current()

	if sxxexxRe.MatchString(token) { // 跳过 SxxExx 格式，让集数解析器处理
		return
	}

	// 检查季数范围格式 "s01-s02"
	if matches := seasonRangeRe.FindStringSubmatch(token + "-" + s.tokens.Peek()); len(matches) > 2 {
		startSeasonStr := matches[1]
		endSeasonStr := matches[2]

		if startSeasonNum, err := strconv.Atoi(startSeasonStr); err == nil && startSeasonNum > 0 {
			if endSeasonNum, err := strconv.Atoi(endSeasonStr); err == nil && endSeasonNum >= startSeasonNum {
				meta.beginSeason = &startSeasonNum
				meta.endSeason = &endSeasonNum
				meta.totalSeason = (endSeasonNum - startSeasonNum) + 1
				meta.mediaType = MediaTypeTV
				s.lastType = lastTokenTypeSeason
				s.stopNameFlag = true
				s.continueFlag = false
				s.tokens.GetNext() // 跳过下一个token（已处理的季数范围）
				return
			}
		}
	}

	// 检查中文季信息格式 "第X季"
	if strings.HasPrefix(token, "第") && strings.HasSuffix(token, "季") {
		// 提取中间的数字
		seasonStr := strings.TrimSuffix(strings.TrimPrefix(token, "第"), "季")

		var seasonNum int = -1
		switch {
		case utils.IsDigits(seasonStr):
			seasonNum, _ = strconv.Atoi(seasonStr)
		case utils.IsAllChinese(seasonStr):
			seasonNum, _ = utils.ChineseToInt(seasonStr)
		}

		if seasonNum > -1 {
			s.lastType = lastTokenTypeSeason
			meta.mediaType = MediaTypeTV
			s.stopcntitleFlag = true // 只停止中文名的处理
			s.continueFlag = false

			if meta.beginSeason == nil {
				meta.beginSeason = &seasonNum
				meta.totalSeason = 1
			}
			return
		}
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
			if matches[i] != "" && utils.IsDigits(matches[i]) {
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
	} else if utils.IsDigits(token) {
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

	// 特殊处理 SxxExx 格式
	if sxxexxMatches := sxxexxRe.FindStringSubmatch(token); len(sxxexxMatches) > 0 {
		seasonStr := sxxexxMatches[1]
		episodeStr := sxxexxMatches[2]

		if seasonNum, err := strconv.Atoi(seasonStr); err == nil && seasonNum > 0 && meta.beginSeason == nil {
			meta.beginSeason = &seasonNum
			meta.totalSeason = 1
		}

		if episodeNum, err := strconv.Atoi(episodeStr); err == nil && episodeNum > 0 && meta.beginEpisode == nil {
			meta.beginEpisode = &episodeNum
			meta.totalEpisode = 1
		}

		goto setEpisodeFlags
	}

	// 处理中文集信息格式
	if strings.HasPrefix(token, "第") {
		var episodeStr string
		switch {
		case strings.HasSuffix(token, "集"):
			episodeStr = strings.TrimSuffix(strings.TrimPrefix(token, "第"), "集")
		case strings.HasSuffix(token, "话"):
			episodeStr = strings.TrimSuffix(strings.TrimPrefix(token, "第"), "话")
		case strings.HasSuffix(token, "話"):
			episodeStr = strings.TrimSuffix(strings.TrimPrefix(token, "第"), "話")
		default:
			goto checkEpisodeRegex
		}

		var episodeNum int
		switch {
		case utils.IsDigits(episodeStr):
			if num, err := strconv.Atoi(episodeStr); err == nil && num > 0 {
				episodeNum = num
			} else {
				goto checkEpisodeRegex
			}
		case utils.IsAllChinese(episodeStr):
			if num, err := utils.ChineseToInt(episodeStr); err == nil && num > 0 {
				episodeNum = num
			} else {
				goto checkEpisodeRegex
			}
		default:
			goto checkEpisodeRegex
		}

		if meta.beginEpisode == nil {
			meta.beginEpisode = &episodeNum
			meta.totalEpisode = 1
		}
		goto setEpisodeFlags
	}

checkEpisodeRegex:
	// 先检查集数范围格式 (如: 01-26)
	if rangeMatches := episodeRangeRe.FindStringSubmatch(token); len(rangeMatches) > 0 {
		startEpisodeStr := rangeMatches[1]
		endEpisodeStr := rangeMatches[2]

		if startEpisode, err := strconv.Atoi(startEpisodeStr); err == nil && startEpisode > 0 {
			if endEpisode, err := strconv.Atoi(endEpisodeStr); err == nil && endEpisode > startEpisode {
				if meta.beginEpisode == nil {
					meta.beginEpisode = &startEpisode
					meta.endEpisode = &endEpisode
					meta.totalEpisode = (endEpisode - startEpisode) + 1
				}
				goto setEpisodeFlags
			}
		}
	}

	// 使用集识别正则匹配
	if matches := episodeRe.FindStringSubmatch(token); len(matches) > 0 {
		var episodeNum int
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" && utils.IsDigits(matches[i]) {
				if num, err := strconv.Atoi(matches[i]); err == nil && num > 0 {
					episodeNum = num
					break
				}
			}
		}

		if episodeNum > 0 {
			if meta.beginEpisode == nil {
				meta.beginEpisode = &episodeNum
				meta.totalEpisode = 1
			} else if episodeNum > *meta.beginEpisode {
				meta.endEpisode = &episodeNum
				meta.totalEpisode = (episodeNum - *meta.beginEpisode) + 1
				if meta.isFile && meta.totalEpisode > 2 {
					meta.endEpisode = nil
					meta.totalEpisode = 1
				}
			}
		}
		goto setEpisodeFlags
	} else if utils.IsDigits(token) { // 处理纯数字token
		tokenInt, _ := strconv.Atoi(token) // 前面已检查过数字，不会出错
		length := len([]rune(token))

		// 检查是否为集数范围格式（当前数字 + 下一个数字）
		nextToken := s.tokens.Peek()
		if utils.IsDigits(nextToken) && len([]rune(nextToken)) >= 1 && len([]rune(nextToken)) <= 4 {
			if nextInt, err := strconv.Atoi(nextToken); err == nil && nextInt > tokenInt && meta.beginEpisode == nil {
				// 这是一个集数范围
				meta.beginEpisode = &tokenInt
				meta.endEpisode = &nextInt
				meta.totalEpisode = (nextInt - tokenInt) + 1
				s.tokens.GetNext() // 跳过下一个 token，因为我们已经处理了
				goto setEpisodeFlags
			}
		}

		switch {
		case meta.beginEpisode != nil && meta.endEpisode == nil &&
			length < 5 && tokenInt > *meta.beginEpisode &&
			s.lastType == lastTokenTypeEpisode:
			meta.endEpisode = &tokenInt
			meta.totalEpisode = (tokenInt - *meta.beginEpisode) + 1
			if meta.isFile && meta.totalEpisode > 2 {
				meta.endEpisode = nil
				meta.totalEpisode = 1
			}
			goto setEpisodeFlags

		case meta.beginEpisode == nil && length > 1 && length < 4 &&
			s.lastType != lastTokenTypeYear && s.lastType != lastTokenTypePix &&
			s.lastType != lastTokenTypeVideoEncode && token != s.unknownNameStr:
			meta.beginEpisode = &tokenInt
			meta.totalEpisode = 1
			goto setEpisodeFlags

		case meta.beginEpisode == nil && length <= 3 && tokenInt <= 99 &&
			s.lastType == lastTokenTypeYear && meta.beginSeason != nil:
			meta.beginEpisode = &tokenInt
			meta.totalEpisode = 1
			goto setEpisodeFlags

		case s.lastType == lastTokenTypeEpisode && meta.beginEpisode == nil && length < 5:
			meta.beginEpisode = &tokenInt
			meta.totalEpisode = 1
			goto setEpisodeFlags
		}
	} else if strings.ToUpper(token) == "EPISODE" { // 处理EPISODE关键词
		s.lastType = lastTokenTypeEpisode
		return
	}

	return

setEpisodeFlags:
	s.lastType = lastTokenTypeEpisode
	s.continueFlag = false
	s.stopNameFlag = true
	meta.mediaType = MediaTypeTV
}

// 识别资源类型
func (meta *MetaVideo) parseResourceType(s *parseState) {
	if meta.GetTitle() == "" {
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
	if meta.GetTitle() == "" {
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
	if meta.GetTitle() == "" {
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

		if meta.videoEncode == encode.VideoEncodeUnknown {
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

			meta.videoEncode = encode.ParseVideoEncode(encodeStr)
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
		if currentIndex >= 1 {
			prevToken = s.tokens.GetByIndex(currentIndex - 1)
		}
		prevTokenUpper := strings.ToUpper(prevToken)

		if prevTokenUpper == "H" || prevTokenUpper == "X" {
			encodeStr := prevTokenUpper + token
			meta.videoEncode = encode.ParseVideoEncode(encodeStr)
			s.continueFlag = false
			s.stopNameFlag = true
		}
		return
	}

	// 处理 VC1、MPEG2 等数字组合
	if utils.IsDigits(token) && s.lastType == lastTokenTypeVideoEncode {
		prevToken := ""
		currentIndex := s.tokens.GetCurrentIndex()
		if currentIndex >= 1 {
			prevToken = s.tokens.GetByIndex(currentIndex - 1)
		}
		prevTokenUpper := strings.ToUpper(prevToken)

		if prevTokenUpper == "VC" || prevTokenUpper == "MPEG" {
			encodeStr := prevTokenUpper + token
			meta.videoEncode = encode.ParseVideoEncode(encodeStr)
			s.continueFlag = false
			s.stopNameFlag = true
		}
		return
	}

	// 处理 10bit 编码
	if tokenUpper == "10BIT" || strings.Contains(tokenUpper, "YUV420P10") {
		s.lastType = lastTokenTypeVideoEncode
		if meta.videoEncode == encode.VideoEncodeUnknown { // 如果没有其他编码信息，设置为纯10bit编码
			meta.videoEncode = encode.VideoEncode10bit
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
	if meta.GetTitle() == "" {
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

	// 贪婪匹配：优先匹配更长的token组合
	separators := []string{"-", " ", ".", ""} // 调整分隔符优先级，DTS-HD MA需要用横杠和空格
	var bestMatch encode.AudioEncode
	var consumedTokens int
	var bestMatchLength int // 记录最佳匹配的长度，用于优先选择更长的匹配

	// 尝试不同长度的token组合（从长到短，贪婪匹配）
	for tokenCount := 1; tokenCount <= 4; tokenCount++ {
		// 构建token列表
		parts := []string{token}

		// 添加后续token
		if tokenCount > 1 {
			for k := 1; k < tokenCount; k++ {
				nextToken := s.tokens.PeekN(k)
				if nextToken == "" {
					goto nextTokenCount // 如果没有更多token，跳到下一个长度
				}
				parts = append(parts, nextToken)
			}
		}

		// 对当前token组合尝试所有分隔符
		for _, separator := range separators {
			var combinedStr string
			if separator == "" {
				combinedStr = strings.Join(parts, "")
			} else if separator == "-" && len(parts) >= 3 && parts[0] == "DTS" {
				// 特殊处理DTS-HD MA格式：DTS-HD MA5.1
				if len(parts) == 3 && parts[1] == "HD" {
					combinedStr = "DTS-HD " + parts[2] // DTS-HD MA5
				} else if len(parts) == 4 && parts[1] == "HD" {
					combinedStr = "DTS-HD " + parts[2] + "." + parts[3] // DTS-HD MA5.1
				} else {
					combinedStr = strings.Join(parts, separator)
				}
			} else {
				combinedStr = strings.Join(parts, separator)
			}

			var testEncode encode.AudioEncode

			// 先使用正则表达式匹配
			if matches := audioEncodeRe.FindStringSubmatch(combinedStr); len(matches) > 0 {
				testEncode = encode.ParseAudioEncode(matches[0])
			} else {
				// 如果正则没匹配到，直接尝试解析
				testEncode = encode.ParseAudioEncode(combinedStr)
			}

			if testEncode != encode.AudioEncodeUnknown {
				// 找到匹配，检查是否比当前最佳匹配更好
				matchLength := len(combinedStr)
				if bestMatch == encode.AudioEncodeUnknown || matchLength > bestMatchLength ||
					(matchLength == bestMatchLength && tokenCount > consumedTokens+1) {
					bestMatch = testEncode
					consumedTokens = tokenCount - 1 // 消费的token数量
					bestMatchLength = matchLength
				}
			}
		}

	nextTokenCount:
		// 如果已经找到匹配且当前长度的组合都尝试完了，可以考虑是否继续
		// 为了真正贪婪匹配，我们继续尝试更长的组合
	}

	// 如果找到了匹配的编码
	if bestMatch != encode.AudioEncodeUnknown {
		s.continueFlag = false
		s.stopNameFlag = true
		s.lastType = lastTokenTypeAudioEncode

		// 消费已使用的token
		for i := 0; i < consumedTokens; i++ {
			s.tokens.GetNext()
		}

		if meta.audioEncode == encode.AudioEncodeUnknown {
			meta.audioEncode = bestMatch
		} else {
			// 如果已有音频编码，进行组合处理
			currentEncode := meta.audioEncode

			// 如果任一编码是Atmos，优先保留Atmos
			if currentEncode == encode.AudioEncodeAtmos || bestMatch == encode.AudioEncodeAtmos {
				meta.audioEncode = encode.AudioEncodeAtmos
			} else {
				// 其他情况使用新匹配的编码（贪婪匹配结果更准确）
				meta.audioEncode = bestMatch
			}
		}
		return
	}

	// 处理数字token（用于音频编码的版本号等）
	if utils.IsDigits(token) && s.lastType == lastTokenTypeAudioEncode {
		if meta.audioEncode != encode.AudioEncodeUnknown {
			currentStr := meta.audioEncode.String()
			if currentStr != "" {
				// 获取前一个token作为参考
				prevToken := ""
				currentIndex := s.tokens.GetCurrentIndex()
				if currentIndex >= 2 {
					prevToken = s.tokens.GetByIndex(currentIndex - 2)
				}

				var newEncodeStr string
				if utils.IsDigits(prevToken) {
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
				newEncode := encode.ParseAudioEncode(newEncodeStr)
				if newEncode != encode.AudioEncodeUnknown {
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
	meta.cntitle = meta.fixName(meta.cntitle)
	meta.entitle = meta.fixName(meta.entitle)

	// 处理BluRay DIY标记
	if meta.resourceType != ResourceTypeUnknown &&
		(meta.resourceType == ResourceTypeBluRay ||
			meta.resourceType == ResourceTypeUHDBluRay ||
			meta.resourceType == ResourceTypeBluRayRemux) {
		// 检查原始字符串中是否包含DIY标记
		upperOriginal := strings.ToUpper(meta.orginalTitle)
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
	if utils.IsDigits(name) {
		num, err := strconv.Atoi(name)
		if err == nil && num < 1800 &&
			meta.year == 0 &&
			meta.beginSeason == nil &&
			meta.resourcePix == ResourcePixUnknown &&
			meta.resourceType == ResourceTypeUnknown &&
			meta.audioEncode == encode.AudioEncodeUnknown &&
			meta.videoEncode == encode.VideoEncodeUnknown {

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
