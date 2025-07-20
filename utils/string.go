package utils

import (
	"MediaTools/extensions"
	"regexp"
	"slices"
	"strings"

	"errors"
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

// IsAllChinese 判断字符串是否全为中文
func IsAllChinese(s string) bool {
	for _, r := range s {
		if r < '\u4e00' || r > '\u9fff' {
			return false
		}
	}
	return len(s) > 0
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

// 中文数字到数值的映射
var chineseNumMap = map[rune]int{
	'零': 0, '一': 1, '二': 2, '三': 3, '四': 4,
	'五': 5, '六': 6, '七': 7, '八': 8, '九': 9,
	'十': 10, '百': 100, '千': 1000, '万': 10000, '亿': 100000000,
}
var ErrNumberTooLarge = errors.New("不支持大于等于一万的数字")
var ErrInvalidCharacter = errors.New("包含无效字符")

// 中文数字转换为整数 (仅支持小于一万的数字)
func ChineseToInt(chinese string) (int, error) {
	if chinese == "" {
		return 0, ErrInvalidCharacter
	}

	// 检查是否包含万或亿，如果包含则返回错误
	if strings.Contains(chinese, "万") || strings.Contains(chinese, "亿") {
		return 0, ErrNumberTooLarge
	}

	// 特殊处理单独的"十"
	if chinese == "十" {
		return 10, nil
	}

	// 如果以"十"开头，前面补"一"
	if strings.HasPrefix(chinese, "十") {
		chinese = "一" + chinese
	}

	result := 0
	current := 0

	for _, char := range chinese {
		if char == '零' {
			continue
		}

		val, exists := chineseNumMap[char]
		if !exists {
			return 0, ErrInvalidCharacter
		}

		if val < 10 { // 数字 0-9
			current = val
		} else if val == 10 { // 十
			if current == 0 {
				current = 1
			}
			current *= 10
			result += current
			current = 0
		} else if val == 100 { // 百
			if current == 0 {
				current = 1
			}
			current *= 100
			result += current
			current = 0
		} else if val == 1000 { // 千
			if current == 0 {
				current = 1
			}
			current *= 1000
			result += current
			current = 0
		}
	}

	// 处理剩余的数字
	if current > 0 {
		// 检查简写情况 - 只有在没有"零"的情况下才视为简写
		hasZero := strings.Contains(chinese, "零")
		if result > 0 && !hasZero {
			// 处理类似"二百五"(250)、"一千二"(1200)的情况
			if result%100 == 0 && result < 1000 {
				// 类似"二百五" -> result=200, current=5, 结果应该是250
				result += current * 10
			} else if result%1000 == 0 && result < 10000 {
				// 类似"一千二" -> result=1000, current=2, 结果应该是1200
				result += current * 100
			} else {
				result += current
			}
		} else {
			result += current
		}
	}

	return result, nil
}
