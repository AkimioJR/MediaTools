package meta

import (
	"regexp"
	"strconv"
)

var (
	// 版本号识别正则
	version1Re     = regexp.MustCompile(`(\d{1,2})[vV](\d+)`)                  // 匹配 03v2 这种格式
	version2Re     = regexp.MustCompile(`\[(\d{1,2}[AaBbVv]\d*)\]`)            // 匹配 [12A] 或 [03v2] 这种格式
	version3Re     = regexp.MustCompile(`(\d{1,2}[AaBb])(?:\s|\]|\[|$|\.|\-)`) // 匹配 12A 或 12B 这种格式，后面跟分隔符
	versionDigitRe = regexp.MustCompile(`[AaBb]|[vV](\d+)`)                        // 提取字母或v后的数字
)

// 解析版本号，返回 uint8 类型版本号（0-255）
// 解析失败时返回 1（默认版本）
func ParseVersion(s string) uint8 {
	// 尝试第1种版本格式 数字+v+数字 (如 03v2)
	if matches := version1Re.FindStringSubmatch(s); len(matches) > 2 {
		if num, err := strconv.ParseUint(matches[2], 10, 8); err == nil {
			return uint8(num)
		}
	}

	// 尝试第2种版本格式 [数字+字母] (如 [12A], [03v2])
	if matches := version2Re.FindStringSubmatch(s); len(matches) > 1 {
		versionStr := matches[1] // 提取括号内的版本部分
		return parseVersionFromString(versionStr)
	}

	// 尝试第3种版本格式 数字+字母 (如 12A, 12B)
	if matches := version3Re.FindStringSubmatch(s); len(matches) > 1 {
		versionStr := matches[1] // 提取版本部分
		return parseVersionFromString(versionStr)
	}

	// 所有格式都不匹配，返回默认版本1
	return 1
}

// 从版本字符串中解析版本号
func parseVersionFromString(versionStr string) uint8 {
	// 查找字母或v+数字的模式
	if matches := versionDigitRe.FindStringSubmatch(versionStr); len(matches) > 0 {
		match := matches[0]
		switch {
		case match == "A" || match == "a":
			return 1
		case match == "B" || match == "b":
			return 2
		case len(matches) > 1 && matches[1] != "": // v后面有数字
			if num, err := strconv.ParseUint(matches[1], 10, 8); err == nil {
				return uint8(num)
			}
		}
	}

	// 如果没有找到版本标识符，返回默认版本1
	return 1
}
