package wordmatch

import (
	"MediaTools/utils"
	"strconv"
	"strings"
)

type WordsMatcher struct {
	words []*CustomWord
}

func NewWordsMatcher(words []string) (*WordsMatcher, error) {
	var matcher WordsMatcher
	matcher.words = make([]*CustomWord, 0, len(words))
	for _, w := range words {
		if strings.HasPrefix(w, "#") || strings.TrimSpace(w) == "" { // 跳过注释和空行
			continue
		}
		cw, err := ParseLine(w)
		if err != nil {
			return nil, err
		}
		matcher.words = append(matcher.words, cw)
	}
	return &matcher, nil
}

func (wm *WordsMatcher) MatchAndProcess(title string) (string, string) {
	var rule string // 匹配到的规则
	for _, word := range wm.words {
		firstMatch := false            // 每次匹配前重置
		if word.replaceFromRe != nil { // 替换被替换词
			if word.replaceFromRe.MatchString(title) {
				rule = word.originalStr
				title = word.replaceFromRe.ReplaceAllString(title, word.ReplaceTo)
				firstMatch = true // 标记为第一次匹配成功
			}
		} else {
			firstMatch = true // 如果没有替换词正则，则直接标记为第一次匹配成功
		}
		if firstMatch && word.PrefixWord != "" && word.SuffixWord != "" && word.OffsetExpr != "" { // 前后定位词和偏移量表达式
			prefixIndex := strings.Index(title, word.PrefixWord)
			suffixIndex := strings.Index(title, word.SuffixWord)

			if prefixIndex == -1 || suffixIndex == -1 || suffixIndex <= prefixIndex {
				continue // 如果没有找到前后定位词，或者后缀在前缀之前，则跳过
			}

			episodeStr := strings.TrimSpace(title[prefixIndex+len(word.PrefixWord) : suffixIndex])
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
				newEpisode, err := ParseOffsetExpr(word.OffsetExpr, episode)
				if err != nil {
					continue
				}
				title = strings.Replace(title, episodeStr, strconv.Itoa(newEpisode), 1)
			}
			rule = word.originalStr
			break
		}
	}
	return title, rule
}
