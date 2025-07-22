package utils

import "regexp"

var removeColorCodesRegexp = regexp.MustCompile(`\033\[[0-9]*m`)

// RemoveColorCodes 移除字符串中的 ANSI 颜色代码
func RemoveColorCodes(line string) string {
	return removeColorCodesRegexp.ReplaceAllString(line, "")
}
