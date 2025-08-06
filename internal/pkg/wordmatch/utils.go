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
