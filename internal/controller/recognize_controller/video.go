package recognize_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

func FormatVideo(item *schemas.MediaItem) (string, error) {
	loock.RLock()
	defer loock.RUnlock()

	var tmpl *template.Template
	switch item.MediaType {
	case meta.MediaTypeMovie:
		tmpl = movieTemplate
	case meta.MediaTypeTV:
		tmpl = tvTemplate
	default:
		return "", fmt.Errorf("不支持的媒体类型: %s", item.MediaType.String())
	}
	var buffer strings.Builder
	if err := tmpl.Execute(&buffer, item); err != nil {
		return "", fmt.Errorf("渲染模板失败: %v", err)
	}
	return buffer.String(), nil
}

// MatchAndProcessVideoTitle 匹配并处理视频标题
// 返回处理后的标题和匹配到的规则
// 如果未匹配到规则，则返回原始标题和空规则
func MatchAndProcessVideoTitle(title string) (string, string) {
	loock.RLock()
	defer loock.RUnlock()

	title, rule := wm.MatchAndProcess(title)
	if rule == "" {
		logrus.Infof("标题 「%s」 未匹配到的规则", title)
	} else {
		logrus.Infof("标题 「%s」 匹配到的规则: 「%s」", title, rule)
	}
	return title, rule
}

func MatchCustomizationWordWord(title string) []string {
	loock.RLock()
	defer loock.RUnlock()

	var customWords = []string{}
	for _, word := range customizationWordRe.FindAllString(title, -1) {
		word = strings.TrimSpace(word)
		if word != "" {
			customWords = append(customWords, word)
		}
	}
	return customWords
}

func ParseVideoMeta(title string) (*meta.VideoMeta, string) {
	loock.RLock()
	defer loock.RUnlock()
	title, rule := MatchAndProcessVideoTitle(title)
	vm := meta.ParseVideoMeta(title)
	vm.Customization = MatchCustomizationWordWord(title)
	return vm, rule
}

var ruleRe = regexp.MustCompile(`\{\[.+\]\}`)

// UpdateMetaByRule 根据匹配规则更新视频元数据
// {[tmdbid=xxx;type=movie/tv;s=xxx;e=xxx]} 直接指定TMDBID，其中s、e为季数和集数（可选）
// 返回应用的规则
func UpdateMetaByRule(vm *meta.VideoMeta) string {
	loock.RLock()
	defer loock.RUnlock()
	matches := ruleRe.FindStringSubmatch(vm.OrginalTitle)
	switch len(matches) {
	case 0:
		logrus.Debugf("标题「%s」未匹配到设置规则", vm.OrginalTitle)
		return ""
	case 1:
		logrus.Debugf("标题「%s」匹配到设置规则：%s", vm.OrginalTitle, matches[0])
		// 解析规则
		rule := matches[0][2 : len(matches[0])-2] // 去掉两边的 {[ 和 ]}
		parts := strings.Split(rule, ";")
		var rules []string
		for _, part := range parts {
			kv := strings.Split(part, "=")
			if len(kv) != 2 {
				logrus.Warningf("标题「%s」匹配到的设置规则格式错误：%s", vm.OrginalTitle, part)
				continue
			}

			switch strings.TrimSpace(strings.ToLower(kv[0])) {
			case "tmdbid": // TMDB ID
				id, err := strconv.Atoi(kv[1])
				if err != nil {
					logrus.Warningf("标题「%s」匹配到的设置规则TMDBID格式错误：%s", vm.OrginalTitle, kv[1])
					continue
				}
				vm.TMDBID = id
				rules = append(rules, "tmdbid="+strconv.Itoa(id))

			case "type": // 媒体类型
				mediaType := meta.ParseMediaType(kv[1])
				if mediaType == meta.MediaTypeUnknown {
					logrus.Warningf("标题「%s」匹配到的设置规则类型错误：%s", vm.OrginalTitle, kv[1])
					continue
				}
				vm.MediaType = mediaType
				rules = append(rules, "type="+mediaType.String())

			case "s": // 季数
				season, err := strconv.Atoi(kv[1])
				if err != nil {
					logrus.Warningf("标题「%s」匹配到的设置规则季数格式错误：%s", vm.OrginalTitle, kv[1])
					continue
				}
				vm.Season = season
				rules = append(rules, "s="+strconv.Itoa(season))

			case "e": // 集数
				episode, err := strconv.Atoi(kv[1])
				if err != nil {
					logrus.Warningf("标题「%s」匹配到的设置规则集数格式错误：%s", vm.OrginalTitle, kv[1])
					continue
				}
				vm.Episode = episode
				rules = append(rules, "e="+strconv.Itoa(episode))
			}
		}
		return "{[" + strings.Join(rules, ";") + "]}"
	default:
		logrus.Warningf("标题「%s」匹配到多个设置规则：%+v，跳过解析", vm.OrginalTitle, matches)
		return ""
	}
}
