package wordmatch

import (
	"MediaTools/utils"
	"strconv"
	"strings"
)

type WordsMatcher struct {
	rules []*CustomWordRule
}

func NewWordsMatcher(lines []string) (*WordsMatcher, error) {
	var matcher WordsMatcher
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" { // 跳过注释和空行
			continue
		}
		cw, err := ParseLine(line)
		if err != nil {
			return nil, err
		}
		matcher.rules = append(matcher.rules, cw)
	}
	return &matcher, nil
}

func (wm *WordsMatcher) MatchAndProcess(title string) (string, string) {
	var rule string // 匹配到的规则
	for _, wordRule := range wm.rules {
		firstMatch := false                // 每次匹配前重置
		if wordRule.replaceFromRe != nil { // 替换被替换词
			if wordRule.replaceFromRe.MatchString(title) {
				rule = wordRule.originalStr
				title = wordRule.replaceFromRe.ReplaceAllString(title, wordRule.ReplaceTo)
				firstMatch = true // 标记为第一次匹配成功
			}
		} else {
			firstMatch = true // 如果没有替换词正则，则直接标记为第一次匹配成功
		}
		if firstMatch && wordRule.PrefixWord != "" && wordRule.SuffixWord != "" && wordRule.OffsetExpr != "" { // 前后定位词和偏移量表达式
			prefixIndex := strings.Index(title, wordRule.PrefixWord)
			suffixIndex := strings.Index(title, wordRule.SuffixWord)

			if prefixIndex == -1 || suffixIndex == -1 || suffixIndex <= prefixIndex {
				continue // 如果没有找到前后定位词，或者后缀在前缀之前，则跳过
			}

			episodeStr := strings.TrimSpace(title[prefixIndex+len(wordRule.PrefixWord) : suffixIndex])
			var episode int
			var err error
			switch {
			case utils.IsDigits(episodeStr):
				episode, err = strconv.Atoi(episodeStr)
				if err != nil {
					continue
				}
			case utils.IsAllChinese(episodeStr):
				episode, err = utils.ChineseToInt(episodeStr)
				if err != nil {
					continue
				}
			case utils.IsRomanNumeral(episodeStr):
				episode, err = utils.RomanToInt(episodeStr)
				if err != nil {
					continue
				}
			default:
				continue
			}
			if episode > 0 {
				newEpisode, err := ParseOffsetExpr(wordRule.OffsetExpr, episode)
				if err != nil {
					continue
				}
				title = strings.Replace(title, episodeStr, strconv.Itoa(newEpisode), 1)
			}
			rule = wordRule.originalStr
			break
		}
	}
	return title, rule
}
