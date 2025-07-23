package meta

import (
	"MediaTools/encode"
	"MediaTools/utils"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// VideoMeta 解析视频媒体名字信息结构体
type VideoMeta struct {
	// 基础信息
	OrginalTitle   string    // 原始标题
	ProcessedTitle string    // 处理后的标题
	IsFile         bool      // 是否是媒体文件
	CNTitle        string    // 中文标题
	ENTitle        string    // 英文标题
	Year           uint      // 年份
	MediaType      MediaType // 媒体类型
	TMDBID         uint64    // TMDB ID

	// 资源信息
	ResourceType   ResourceType                // 来源/介质
	ResourceEffect map[ResourceEffect]struct{} // 资源效果
	ResourcePix    ResourcePix                 // 分辨率
	VideoEncode    encode.VideoEncode          // 视频编码
	AudioEncode    encode.AudioEncode          // 音频编码
	Platform       StreamingPlatform           // 流媒体平台
	ReleaseGroups  []string                    // 发布组
	Part           string                      // 分段
	Version        uint8                       // 版本号
	// customization  string                      // 自定义词

	// 电视剧相关·
	BeginSeason  *int // 起始季
	EndSeason    *int // 结束集
	TotalSeason  int  // 总季数
	BeginEpisode *int // 起始集
	EndEpisode   *int // 结束集
	TotalEpisode int  // 总集数
}

// 获取标题
// 有中文标题优先返回中文标题
// 否则返回英文标题
func (meta *VideoMeta) GetTitle() string {
	if meta.CNTitle != "" {
		return meta.CNTitle
	}
	return meta.ENTitle
}

// 获取标题列表
// 返回中文标题和英文标题的列表
func (meta *VideoMeta) GetTitles() []string {
	titles := make([]string, 0, 2)
	if meta.CNTitle != "" {
		titles = append(titles, meta.CNTitle)
	}
	if meta.ENTitle != "" {
		titles = append(titles, meta.ENTitle)
	}
	return titles
}

func (meta *VideoMeta) GetResourceEffectStrings() []string {
	var effects []string
	for effect := range meta.ResourceEffect {
		effects = append(effects, effect.String())
	}
	return effects
}

func (meta *VideoMeta) GetSeasons() []int {
	if meta.MediaType != MediaTypeTV {
		return nil
	}

	var seasons []int
	if meta.BeginSeason != nil {
		if meta.EndSeason != nil {
			for i := *meta.BeginSeason; i <= *meta.EndSeason; i++ {
				seasons = append(seasons, i)
			}
		} else {
			seasons = append(seasons, *meta.BeginSeason)
		}
	}
	return seasons
}

func (meta *VideoMeta) GetSeasonStr() string {
	if meta.BeginSeason == nil {
		return ""
	} else {
		if meta.EndSeason == nil {
			return fmt.Sprintf("S%02d", *meta.BeginSeason)
		}
		return fmt.Sprintf("S%02d-S%02d", *meta.BeginSeason, *meta.EndSeason)
	}
}

// GetEpisodes 获取集数列表
func (meta *VideoMeta) GetEpisodes() []int {
	if meta.MediaType != MediaTypeTV {
		return nil
	}

	var episodes []int
	if meta.BeginEpisode != nil {
		if meta.EndEpisode != nil {
			for i := *meta.BeginEpisode; i <= *meta.EndEpisode; i++ {
				episodes = append(episodes, i)
			}
		} else {
			episodes = append(episodes, *meta.BeginEpisode)
		}
	}
	return episodes
}

func (meta *VideoMeta) GetEpisodeStr() string {
	if meta.BeginEpisode == nil {
		return ""
	} else {
		if meta.EndEpisode == nil {
			return fmt.Sprintf("E%02d", *meta.BeginEpisode)
		}
		return fmt.Sprintf("E%02d-E%02d", *meta.BeginEpisode, *meta.EndEpisode)
	}
}

func ParseVideoMeta(title string) *VideoMeta {
	meta := &VideoMeta{
		OrginalTitle:   title,
		MediaType:      MediaTypeUnknown,
		ResourceType:   ResourceTypeUnknown,
		ResourceEffect: make(map[ResourceEffect]struct{}),
		ReleaseGroups:  findReleaseGroups(title), // 解析发布组
		Platform:       UnknownStreamingPlatform,
		Version:        ParseVersion(title), // 解析版本号
	}

	if utils.IsMediaExtension(path.Ext(title)) {
		title = strings.TrimSuffix(title, path.Ext(title)) // 去掉文件扩展名
		meta.IsFile = true
	}

	loc := nameNoBeginRe.FindStringIndex(title) // 去掉名称中第1个[]的内容（一般是发布组）
	if loc != nil {
		title = title[:loc[0]] + title[loc[1]:]
	}
	title = yearRangeRe.ReplaceAllString(title, "${1}${2}") // 把xxxx-xxxx年份换成前一个年份，常出现在季集上
	title = fileSizeRe.ReplaceAllString(title, "")          // 把大小去掉
	title = dateFmtRe.ReplaceAllString(title, "")           // 把年月日去掉
	title = strings.TrimSpace(title)                        // 去掉首尾空格
	meta.ProcessedTitle = title

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
func (meta *VideoMeta) parsePart(s *parseState) {
	if meta.GetTitle() == "" {
		return
	}
	if meta.Year == 0 &&
		meta.BeginSeason == nil &&
		meta.BeginEpisode == nil &&
		meta.ResourcePix == ResourcePixUnknown &&
		meta.ResourceType == ResourceTypeUnknown {
		return
	}

	token := s.tokens.Current()

	if partRe.MatchString(token) {
		meta.Part = token

		nextToken := s.tokens.Peek()
		utf8Str := []rune(nextToken)
		length := len(utf8Str)
		if nextToken != "" {
			if (utils.IsDigits(nextToken) && (length == 1 || (length == 2 && utf8Str[0] == '0'))) ||
				slices.Contains([]string{"A", "B", "C", "I", "II", "III"}, strings.ToUpper(nextToken)) {
				meta.Part += nextToken
			}
		}

		s.lastType = lastTokenTypePart
		s.continueFlag = false
	}
}

// 识别 cntitle、entitle
func (meta *VideoMeta) parseName(s *parseState) {
	token := s.tokens.Current()

	if s.unknownNameStr != "" { // 回收标题
		if meta.CNTitle == "" {
			if meta.ENTitle == "" {
				meta.ENTitle = s.unknownNameStr
			} else if s.unknownNameStr != strconv.Itoa(int(meta.Year)) {
				meta.ENTitle += " " + s.unknownNameStr
			}
			s.lastType = lastTokenTypeentitle
		}
		s.unknownNameStr = ""
	}

	if s.stopNameFlag {
		return
	}

	if slices.Contains(meta.ReleaseGroups, token) { // 如果当前token是发布组，直接跳过
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
		if meta.CNTitle == "" {
			meta.CNTitle = token
		} else if !s.stopcntitleFlag {
			// 含有电影关键词或者不含特殊字符的中文可以继续拼接
			if slices.Contains([]string{"剧场版", "劇場版", "电影版", "電影版"}, token) ||
				(!nameNoChineseRe.MatchString(token) && !slices.Contains([]string{"共", "第", "季", "集", "话", "話", "期"}, token)) {
				meta.CNTitle += " " + token
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

					if !utils.IsRomanNumeral(token) && // 中文名后面跟的数字不是年份的极有可能是集
						s.lastType == lastTokenTypecntitle &&
						tokenInt < YearMin {
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
								switch s.lastType {
								case lastTokenTypecntitle:
									currentTitle = meta.CNTitle
								case lastTokenTypeentitle:
									currentTitle = meta.ENTitle
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
										meta.CNTitle += " " + token
									case lastTokenTypeentitle:
										meta.ENTitle += " " + token
									}
									s.continueFlag = false
									return
								}

								// 如果数字过大（如 >100），不太可能是简单的集数，更可能是标题的一部分
								if tokenInt > 100 {
									switch s.lastType {
									case lastTokenTypecntitle:
										meta.CNTitle += " " + token
									case lastTokenTypeentitle:
										meta.ENTitle += " " + token
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
									meta.CNTitle += " " + token
								case lastTokenTypeentitle:
									meta.ENTitle += " " + token
								}
								s.continueFlag = false
								return
							}
						}
					}
					// 其他情况（罗马数字或其他特殊数字）拼装到已有标题中
					switch s.lastType {
					case lastTokenTypecntitle:
						meta.CNTitle += " " + token
					case lastTokenTypeentitle:
						meta.ENTitle += " " + token
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
			if meta.ENTitle != "" && strings.HasSuffix(strings.ToUpper(meta.ENTitle), "SEASON") { // 如果匹配到季，英文名结尾为Season，说明Season属于标题，不应在后续作为干扰词去除
				meta.ENTitle += " "
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
			if meta.ENTitle != "" {
				meta.ENTitle += " " + token
			} else {
				meta.ENTitle = token
			}
			s.lastType = lastTokenTypeentitle
		}
	}
}

// 识别年份
func (meta *VideoMeta) parseYear(s *parseState) {
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
	if num < YearMin || num > YearMax {
		return
	}

	if meta.Year != 0 {
		if meta.ENTitle != "" {
			meta.ENTitle = strings.TrimSpace(meta.ENTitle) + " " + strconv.Itoa(int(meta.Year))
		} else if meta.CNTitle != "" {
			meta.CNTitle += " " + strconv.Itoa(int(meta.Year))
		}
	} else if meta.ENTitle != "" && strings.HasSuffix(strings.ToLower(meta.ENTitle), "season") { // 如果匹配到年，且英文名结尾为Season，说明Season属于标题，不应在后续作为干扰词去除
		meta.ENTitle += " "
	}

	meta.Year = uint(num)
	s.lastType = lastTokenTypeYear
	s.continueFlag = false
	s.stopNameFlag = true
}

// 识别分辨率
func (meta *VideoMeta) parseResourcePix(s *parseState) {
	if meta.GetTitle() == "" {
		return
	}

	token := s.tokens.Current()
	var r ResourcePix = ResourcePixUnknown

	// 先使用更精确的2K/4K/8K匹配
	matches := resourcePixStandardRe.FindStringSubmatch(token)
	if len(matches) > 1 {
		if meta.ResourcePix == ResourcePixUnknown {
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
				if resourcePixStr != "" && meta.ResourcePix == ResourcePixUnknown {
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
	meta.ResourcePix = r
}

// 识别季
func (meta *VideoMeta) parseSeason(s *parseState) {
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
				meta.BeginSeason = &startSeasonNum
				meta.EndSeason = &endSeasonNum
				meta.TotalSeason = (endSeasonNum - startSeasonNum) + 1
				meta.MediaType = MediaTypeTV
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
			meta.MediaType = MediaTypeTV
			s.stopcntitleFlag = true // 只停止中文名的处理
			s.continueFlag = false

			if meta.BeginSeason == nil {
				meta.BeginSeason = &seasonNum
				meta.TotalSeason = 1
			}
			return
		}
	}

	// 使用季识别正则匹配
	matches := seasonRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		s.lastType = lastTokenTypeSeason
		meta.MediaType = MediaTypeTV
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
			if meta.BeginSeason == nil {
				meta.BeginSeason = &seasonNum
				meta.TotalSeason = 1
			} else {
				if seasonNum > *meta.BeginSeason {
					meta.EndSeason = &seasonNum
					meta.TotalSeason = (seasonNum - *meta.BeginSeason) + 1
					// 如果是文件且总季数大于1，重置结束季
					if meta.IsFile && meta.TotalSeason > 1 {
						meta.EndSeason = nil
						meta.TotalSeason = 1
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
			meta.BeginSeason == nil &&
			len(token) < 3 {
			meta.BeginSeason = &tokenInt
			meta.TotalSeason = 1
			s.lastType = lastTokenTypeSeason
			s.stopNameFlag = true
			s.continueFlag = false
			meta.MediaType = MediaTypeTV
		}
	} else if strings.ToUpper(token) == "SEASON" && meta.BeginSeason == nil {
		// 遇到SEASON关键词
		s.lastType = lastTokenTypeSeason
	} else if meta.MediaType == MediaTypeTV && meta.BeginSeason == nil {
		// 如果已确定为电视剧类型但没有季数，默认为第1季
		defaultSeason := 1
		meta.BeginSeason = &defaultSeason
	}
}

// 识别集数
func (meta *VideoMeta) parseEpisode(s *parseState) {
	token := s.tokens.Current()

	// 特殊处理 SxxExx 格式
	if sxxexxMatches := sxxexxRe.FindStringSubmatch(token); len(sxxexxMatches) > 0 {
		seasonStr := sxxexxMatches[1]
		episodeStr := sxxexxMatches[2]

		if seasonNum, err := strconv.Atoi(seasonStr); err == nil && seasonNum > 0 && meta.BeginSeason == nil {
			meta.BeginSeason = &seasonNum
			meta.TotalSeason = 1
		}

		if episodeNum, err := strconv.Atoi(episodeStr); err == nil && episodeNum > 0 && meta.BeginEpisode == nil {
			meta.BeginEpisode = &episodeNum
			meta.TotalEpisode = 1
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

		if meta.BeginEpisode == nil {
			meta.BeginEpisode = &episodeNum
			meta.TotalEpisode = 1
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
				if meta.BeginEpisode == nil {
					meta.BeginEpisode = &startEpisode
					meta.EndEpisode = &endEpisode
					meta.TotalEpisode = (endEpisode - startEpisode) + 1
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
			if meta.BeginEpisode == nil {
				meta.BeginEpisode = &episodeNum
				meta.TotalEpisode = 1
			} else if episodeNum > *meta.BeginEpisode {
				meta.EndEpisode = &episodeNum
				meta.TotalEpisode = (episodeNum - *meta.BeginEpisode) + 1
				if meta.IsFile && meta.TotalEpisode > 2 {
					meta.EndEpisode = nil
					meta.TotalEpisode = 1
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
			if nextInt, err := strconv.Atoi(nextToken); err == nil && nextInt > tokenInt && meta.BeginEpisode == nil {
				// 这是一个集数范围
				meta.BeginEpisode = &tokenInt
				meta.EndEpisode = &nextInt
				meta.TotalEpisode = (nextInt - tokenInt) + 1
				s.tokens.GetNext() // 跳过下一个 token，因为我们已经处理了
				goto setEpisodeFlags
			}
		}

		switch {
		case meta.BeginEpisode != nil && meta.EndEpisode == nil &&
			length < 5 && tokenInt > *meta.BeginEpisode &&
			s.lastType == lastTokenTypeEpisode:
			meta.EndEpisode = &tokenInt
			meta.TotalEpisode = (tokenInt - *meta.BeginEpisode) + 1
			if meta.IsFile && meta.TotalEpisode > 2 {
				meta.EndEpisode = nil
				meta.TotalEpisode = 1
			}
			goto setEpisodeFlags

		case meta.BeginEpisode == nil && length > 1 && length < 4 &&
			s.lastType != lastTokenTypeYear && s.lastType != lastTokenTypePix &&
			s.lastType != lastTokenTypeVideoEncode && token != s.unknownNameStr:
			meta.BeginEpisode = &tokenInt
			meta.TotalEpisode = 1
			goto setEpisodeFlags

		case meta.BeginEpisode == nil && length <= 3 && tokenInt <= 99 &&
			s.lastType == lastTokenTypeYear && meta.BeginSeason != nil:
			meta.BeginEpisode = &tokenInt
			meta.TotalEpisode = 1
			goto setEpisodeFlags

		case s.lastType == lastTokenTypeEpisode && meta.BeginEpisode == nil && length < 5:
			meta.BeginEpisode = &tokenInt
			meta.TotalEpisode = 1
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
	meta.MediaType = MediaTypeTV
}

// 识别资源类型
func (meta *VideoMeta) parseResourceType(s *parseState) {
	if meta.GetTitle() == "" {
		return
	}

	token := s.tokens.Current()
	tokenUpper := strings.ToUpper(token)

	// 处理特殊组合情况
	if tokenUpper == "DL" && s.lastType == lastTokenTypeSource && meta.ResourceType == ResourceTypeWeb {
		meta.ResourceType = ResourceTypeWebDL
		s.continueFlag = false
		return
	} else if tokenUpper == "RAY" && s.lastType == lastTokenTypeSource && meta.ResourceType == ResourceTypeBlu {
		// UHD BluRay组合
		if meta.ResourceType == ResourceTypeUHD {
			meta.ResourceType = ResourceTypeUHDBluRay
		} else {
			meta.ResourceType = ResourceTypeBluRay
		}
		s.continueFlag = false
		return
	} else if tokenUpper == "WEBDL" {
		meta.ResourceType = ResourceTypeWebDL
		s.continueFlag = false
		return
	}

	// UHD REMUX组合
	if tokenUpper == "REMUX" && meta.ResourceType == ResourceTypeBluRay {
		meta.ResourceType = ResourceTypeBluRayRemux
		s.continueFlag = false
		return
	} else if tokenUpper == "BLURAY" && meta.ResourceType == ResourceTypeUHD {
		meta.ResourceType = ResourceTypeUHDBluRay
		s.continueFlag = false
		return
	}

	// 使用资源类型正则匹配
	matches := sourceRe.FindStringSubmatch(token)
	if len(matches) > 0 {
		s.lastType = lastTokenTypeSource
		s.continueFlag = false
		s.stopNameFlag = true
		if meta.ResourceType == ResourceTypeUnknown {
			meta.ResourceType = ParseResourceType(matches[0])
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
			meta.ResourceEffect[effect] = struct{}{}
		}
		return
	}
}

// 识别流媒体平台
func (meta *VideoMeta) parsePlatform(s *parseState) {
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
		meta.Platform = ParseStreamingPlatform(platformName)
		s.lastType = lastTokenTypePlatform
		s.continueFlag = false
		s.stopNameFlag = true
	}
}

// 识别视频编码
func (meta *VideoMeta) parseVideoEncode(s *parseState) {
	// 检查是否已有名称
	if meta.GetTitle() == "" {
		return
	}

	// 检查是否有其他必要信息
	if meta.Year == 0 &&
		meta.ResourcePix == ResourcePixUnknown &&
		meta.ResourceType == ResourceTypeUnknown &&
		meta.BeginSeason == nil &&
		meta.BeginEpisode == nil {
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

		if meta.VideoEncode == encode.VideoEncodeUnknown {
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

			meta.VideoEncode = encode.ParseVideoEncode(encodeStr)
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
			meta.VideoEncode = encode.ParseVideoEncode(encodeStr)
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
			meta.VideoEncode = encode.ParseVideoEncode(encodeStr)
			s.continueFlag = false
			s.stopNameFlag = true
		}
		return
	}

	// 处理 10bit 编码
	if tokenUpper == "10BIT" || strings.Contains(tokenUpper, "YUV420P10") {
		s.lastType = lastTokenTypeVideoEncode
		if meta.VideoEncode == encode.VideoEncodeUnknown { // 如果没有其他编码信息，设置为纯10bit编码
			meta.VideoEncode = encode.VideoEncode10bit
		} else { // 使用辅助方法升级为对应的10bit版本
			meta.VideoEncode = meta.VideoEncode.CombineWith10bit()
		}
		s.continueFlag = false
		s.stopNameFlag = true
		return
	}
}

// 识别音频编码
func (meta *VideoMeta) parseAudioEncode(s *parseState) {
	// 检查是否已有名称
	if meta.GetTitle() == "" {
		return
	}

	// 检查是否有其他必要信息
	if meta.Year == 0 &&
		meta.ResourcePix == ResourcePixUnknown &&
		meta.ResourceType == ResourceTypeUnknown &&
		meta.BeginSeason == nil &&
		meta.BeginEpisode == nil {
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

		if meta.AudioEncode == encode.AudioEncodeUnknown {
			meta.AudioEncode = bestMatch
		} else {
			// 如果已有音频编码，进行组合处理
			currentEncode := meta.AudioEncode

			// 如果任一编码是Atmos，优先保留Atmos
			if currentEncode == encode.AudioEncodeAtmos || bestMatch == encode.AudioEncodeAtmos {
				meta.AudioEncode = encode.AudioEncodeAtmos
			} else {
				// 其他情况使用新匹配的编码（贪婪匹配结果更准确）
				meta.AudioEncode = bestMatch
			}
		}
		return
	}

	// 处理数字token（用于音频编码的版本号等）
	if utils.IsDigits(token) && s.lastType == lastTokenTypeAudioEncode {
		if meta.AudioEncode != encode.AudioEncodeUnknown {
			currentStr := meta.AudioEncode.String()
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
					meta.AudioEncode = newEncode
				}
				// 注：如果解析失败，保持原有编码不变
			}
		}
		s.continueFlag = false
		return
	}
}

// postProcess 后处理：清理和优化解析结果
func (meta *VideoMeta) postProcess() {
	// 处理part
	if meta.Part != "" && strings.ToUpper(meta.Part) == "PART" {
		meta.Part = ""
	}

	// 清理名称中的干扰字符
	meta.CNTitle = meta.fixName(meta.CNTitle)
	meta.ENTitle = meta.fixName(meta.ENTitle)

	// 处理BluRay DIY标记
	if meta.ResourceType != ResourceTypeUnknown &&
		(meta.ResourceType == ResourceTypeBluRay ||
			meta.ResourceType == ResourceTypeUHDBluRay ||
			meta.ResourceType == ResourceTypeBluRayRemux) {
		// 检查原始字符串中是否包含DIY标记
		upperOriginal := strings.ToUpper(meta.OrginalTitle)
		if strings.Contains(upperOriginal, "DIY") ||
			strings.Contains(upperOriginal, "-DIY@") {
			// 可以添加DIY标记到资源效果中
			meta.ResourceEffect[ResourceEffectUnknown] = struct{}{} // 需要定义DIY效果类型
		}
	}
}

// fixName 清理名称中的干扰字符
func (meta *VideoMeta) fixName(name string) string {
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
			meta.Year == 0 &&
			meta.BeginSeason == nil &&
			meta.ResourcePix == ResourcePixUnknown &&
			meta.ResourceType == ResourceTypeUnknown &&
			meta.AudioEncode == encode.AudioEncodeUnknown &&
			meta.VideoEncode == encode.VideoEncodeUnknown {

			// 如果还没有起始集，将此数字设为起始集
			if meta.BeginEpisode == nil {
				meta.BeginEpisode = &num
				meta.TotalEpisode = 1
				meta.MediaType = MediaTypeTV
				return ""
			} else if meta.isInEpisode(num) && meta.BeginSeason == nil {
				// 如果数字在集数范围内且没有季数信息，清空名称
				return ""
			}
		}
	}

	return name
}

// isInEpisode 检查数字是否在当前集数范围内
func (meta *VideoMeta) isInEpisode(episode int) bool {
	if meta.BeginEpisode == nil {
		return false
	}

	if meta.EndEpisode == nil {
		return episode == *meta.BeginEpisode
	}

	return episode >= *meta.BeginEpisode && episode <= *meta.EndEpisode
}

func ParseVideoMetaByPath(p string) *VideoMeta {
	idx := 3
	parts := strings.Split(p, "/")
	names := make([]string, idx)
	for i := len(parts) - 1; i >= 0; i-- {
		name := strings.TrimSpace(nameMovieWordsRe.ReplaceAllString(parts[i], ""))
		if name != "" {
			idx--
			names[idx] = name
		}
		if idx == 0 { // 只取最后三个部分作为名称
			break
		}
	}
	return ParseVideoMeta(strings.Join(names, " "))
}
