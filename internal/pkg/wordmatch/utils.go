package wordmatch

import (
	"strconv"
	"strings"
)

func ParseOffsetExpr(expr string, episode int) (int, error) {
	if expr == "" {
		return episode, nil
	}
	strs := strings.Split(expr, "EP")
	if len(strs) != 2 {
		return 0, ErrInvalidEpisodeOffsetFormat
	}
	strs[0] = strings.TrimSpace(strings.ReplaceAll(strs[0], "*", ""))
	strs[1] = strings.TrimSpace(strs[1])

	var timesOffset int = 1
	var offset = 0
	var err error
	if strs[0] != "" {
		timesOffset, err = strconv.Atoi(strs[0])
		if err != nil {
			return 0, ErrInvalidEpisodeOffsetFormat
		}
	}
	if strs[1] != "" {
		offset, err = strconv.Atoi(strs[1])
		if err != nil {
			return 0, ErrInvalidEpisodeOffsetFormat
		}
	}
	return timesOffset*episode + offset, nil
}

func ParseReplaceExpr(expr string) ([2]string, error) {
	strs := strings.SplitN(expr, " => ", 2)
	if len(strs) != 2 {
		return [2]string{"", ""}, ErrInvalidReplaceExprFormat
	}
	return [2]string{strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1])}, nil
}

func ParseEpisodeExpr(expr string) ([3]string, error) {
	strs1 := strings.SplitN(expr, " <> ", 2)
	if len(strs1) != 2 {
		return [3]string{"", "", ""}, ErrInvalidEpisodeExprFormat
	}
	strs2 := strings.SplitN(strs1[1], " >> ", 2)
	if len(strs2) != 2 {
		return [3]string{"", "", ""}, ErrInvalidEpisodeExprFormat
	}
	return [3]string{strings.TrimSpace(strs1[0]), strings.TrimSpace(strs2[0]), strings.TrimSpace(strs2[1])}, nil
}

func ParseLine(s string) (*CustomWord, error) {
	var (
		replceStrs  [2]string
		episodeStrs [3]string
		err         error
	)
	hasArrow := strings.Contains(s, " => ")
	hasEpisodeExpr := strings.Contains(s, " <> ") && strings.Contains(s, " >> ")
	strs := strings.Split(s, " && ")
	switch {
	case hasArrow && hasEpisodeExpr: // 有替换表达式和前后定位词偏移量表达式
		// 被替换词 => 替换词 && 前定位词 <> 后定位词 >> 集偏移量（EP）
		strs[0] = strings.TrimSpace(strs[0])
		strs[1] = strings.TrimSpace(strs[1])

		if len(strs) != 2 || strs[0] == "" || strs[1] == "" {
			return nil, ErrInvalidLineFormat
		}

		if strs[0] != "" {
			replceStrs, err = ParseReplaceExpr(strs[0])
			if err != nil {
				return nil, err
			}
		}
		if strs[1] != "" {
			episodeStrs, err = ParseEpisodeExpr(strs[1])
			if err != nil {
				return nil, err
			}
		}
	case hasArrow && !hasEpisodeExpr: // 只有替换表达式
		// 被替换词 => 替换词
		strs[0] = strings.TrimSpace(strs[0])
		if strs[0] == "" {
			return nil, ErrInvalidLineFormat
		}
		replceStrs, err = ParseReplaceExpr(strs[0])
		if err != nil {
			return nil, err
		}
	case hasEpisodeExpr && !hasArrow: // 只有前后定位词偏移量表达式
		// 前定位词 <> 后定位词 >> 集偏移量（EP）
		strs[0] = strings.TrimSpace(strs[0])
		if strs[0] == "" {
			return nil, ErrInvalidLineFormat
		}
		episodeStrs, err = ParseEpisodeExpr(strs[0])
		if err != nil {
			return nil, err
		}
	case !hasArrow && !hasEpisodeExpr: // 纯屏蔽词
		// 被屏蔽词
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			return nil, ErrInvalidLineFormat
		}
		replceStrs[0] = trimmed
	}

	cw := &CustomWord{
		ReplaceFrom: replceStrs[0],
		ReplaceTo:   replceStrs[1],
		PrefixWord:  episodeStrs[0],
		SuffixWord:  episodeStrs[1],
		OffsetExpr:  episodeStrs[2],
	}

	if err := cw.Compile(); err != nil {
		return nil, err
	}

	return cw, nil
}
