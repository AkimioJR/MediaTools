package meta

import (
	"strings"
)

// isDigits 判断字符串是否全为数字
func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// isChinese 判断字符串是否包含中文
func isChinese(s string) bool {
	for _, r := range s {
		if r >= '\u4e00' && r <= '\u9fff' {
			return true
		}
	}
	return false
}

func contain[T comparable](arr []T, ele T) bool {
	for _, e := range arr {
		if ele == e {
			return true
		}
	}
	return false
}

// isRomanNumeral 判断是否为罗马数字
func isRomanNumeral(s string) bool {
	return romanNumeralsRe.MatchString(strings.ToUpper(s))
}

// isMediaExtension 判断是否为媒体文件扩展名
func isMediaExtension(s string) bool {
	return contain(MediaExtensions, strings.ToLower(s))
}
