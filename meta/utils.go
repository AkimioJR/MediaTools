package meta

import (
	"regexp"
	"strconv"
)

var (
	// 版本号识别正则
	version1Re = regexp.MustCompile(`(\d{2})[vV](\d+)`)              // 匹配 03v2 这种格式
	version2Re = regexp.MustCompile(`\[(\d{0,2}([AaBb]|[Vv]\d)?)\]`) // 匹配 [03v2] 或 [1A] 这种格式
	version3Re = regexp.MustCompile(` \d{1,2}[AaBb] `)               // 匹配 1A 或 2B 这种格式
)

// 解析版本号，返回 uint8 类型版本号（0-255）
// 解析失败时返回 0
func ParseVersion(s string) uint8 {
	// 尝试第1种版本格式 数字+v+数字 (如 03v2)
	if matches := version1Re.FindStringSubmatch(s); len(matches) > 2 {
		if num, err := strconv.ParseUint(matches[2], 10, 8); err == nil {
			return uint8(num)
		}
	}

	// 尝试第2种版本格式 [数字][字母]
	if matches := version2Re.FindStringSubmatch(s); len(matches) > 1 {
		versionStr := matches[1] // 提取括号内的版本部分

		// 检查是否包含字母，如果包含则返回对应的版本号
		for i, c := range versionStr {
			switch c {
			case 'A', 'a':
				return 1
			case 'B', 'b':
				return 2
			case 'V', 'v':
				// 如果是V后面跟数字，提取V后面的数字
				if i+1 < len(versionStr) {
					numStr := versionStr[i+1:]
					if num, err := strconv.ParseUint(numStr, 10, 8); err == nil {
						return uint8(num)
					}
				}
			}
		}

		// 如果没有字母，尝试提取纯数字部分
		numStr := ""
		for _, c := range versionStr {
			if c >= '0' && c <= '9' {
				numStr += string(c)
			}
		}

		// 转换数字部分
		if numStr != "" {
			if num, err := strconv.ParseUint(numStr, 10, 8); err == nil {
				return uint8(num)
			}
		}
	}

	// 尝试第3种版本格式 数字+字母+空格
	if matches := version3Re.FindString(s); matches != "" {
		// 检查字母类型
		for _, c := range matches {
			switch c {
			case 'A', 'a':
				return 1
			case 'B', 'b':
				return 2
			}
		}
	}

	// 所有格式都不匹配
	return 0
}
