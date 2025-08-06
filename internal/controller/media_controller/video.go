package media_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"
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

	matchedWords := customizationWordRe.FindAllString(title, -1)
	if len(matchedWords) == 0 {
		return []string{}
	}
	return matchedWords
}

func ParseVideoMeta(title string) *meta.VideoMeta {
	loock.RLock()
	defer loock.RUnlock()
	title, _ = MatchAndProcessVideoTitle(title)
	vm := meta.ParseVideoMeta(title)
	vm.Customization = MatchCustomizationWordWord(title)
	return vm
}
