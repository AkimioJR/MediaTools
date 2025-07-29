package fanart

import (
	"slices"
	"sort"
)

// 根据 languages 优先级排序 BaseInfo
func SortByLanguages[T img](infos []T, languages []string) {
	sort.Slice(infos, func(i, j int) bool {
		return slices.Index(languages, infos[i].getLang()) < slices.Index(languages, infos[j].getLang())
	})
}
