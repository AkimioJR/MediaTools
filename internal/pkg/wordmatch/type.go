package wordmatch

import (
	"errors"
	"regexp"
)

type CustomWord struct {
	ReplaceFrom string `json:"replace_from"` // 被替换词
	ReplaceTo   string `json:"replace_to"`   // 替换词
	PrefixWord  string `json:"prefix_word"`  // 前定位词
	SuffixWord  string `json:"suffix_word"`  // 后定位词
	OffsetExpr  string `json:"offset_expr"`  // 偏移量表达式

	replaceFromRe *regexp.Regexp // 被替换词正则
}

func (cw *CustomWord) Compile() error {
	var err error
	if cw.ReplaceFrom != "" {
		cw.replaceFromRe, err = regexp.Compile(cw.ReplaceFrom)
		if err != nil {
			return err
		}
	}
	return nil
}

var ErrInvalidLineFormat = errors.New("invalid line format")
var ErrInvalidEpisodeOffsetFormat = errors.New("invalid episode offset format")
var ErrInvalidReplaceExprFormat = errors.New("invalid replace expression format")
var ErrInvalidEpisodeExprFormat = errors.New("invalid episode expression format")
