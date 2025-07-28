package fanart

import (
	"slices"
	"sort"
)

// 根据 languages 优先级排序并过滤 images，只有 Lang 在 languages 中的才保留
func SortByLanguages[T img](images []T, languages []string) []T {
	result := make([]T, 0, len(images))

	// 过滤不在 langs 中的图片
	for _, img := range images {
		if !slices.Contains(languages, img.getLang()) {
			continue
		}
		result = append(result, img)
	}

	// 根据 langs 的顺序排序
	sort.Slice(result, func(i, j int) bool {
		return slices.Index(languages, result[i].getLang()) < slices.Index(languages, result[j].getLang())
	})

	return result
}
