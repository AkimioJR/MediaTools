package utils

import (
	"MediaTools/extensions"
	"regexp"
	"slices"
	"strings"
)

// IsDigits 判断字符串是否全为数字
func IsDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// IsChinese 判断字符串是否包含中文
func IsChinese(s string) bool {
	for _, r := range s {
		if r >= '\u4e00' && r <= '\u9fff' {
			return true
		}
	}
	return false
}

var romanNumeralsRe = regexp.MustCompile(`^M*(?:C[MD]|D?C{0,3})(?:X[CL]|L?X{0,3})(?:I[XV]|V?I{0,3})$`) // 罗马数字识别
// IsRomanNumeral 判断是否为罗马数字
func IsRomanNumeral(s string) bool {
	return romanNumeralsRe.MatchString(strings.ToUpper(s))
}

// IsMediaExtension 判断是否为媒体文件扩展名
func IsMediaExtension(s string) bool {
	return slices.Contains(extensions.MediaExtensions, strings.ToLower(s))
}
