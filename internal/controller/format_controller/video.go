package format_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"
	"strings"
	"text/template"
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
