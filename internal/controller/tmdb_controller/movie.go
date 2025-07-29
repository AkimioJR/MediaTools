package tmdb_controller

import (
	"MediaTools/internal/pkg/themoviedb/v3"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// 根据电影的名称和年份匹配TMDB
// name: 电影名称
// year: 电影年份，传入 nil 则不限制年份
func SearchMovieByName(name string, year *int) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()
	
	var params themoviedb.SearchMovieParams
	params.Query = name
	if year != nil && *year > 0 {
		yearStr := strconv.Itoa(*year)
		params.Year = &yearStr
		logrus.Infof("正在搜索电影「%s (%d)」...", name, *year)
	} else {
		logrus.Infof("正在搜索电影「%s」...", name)
	}

	resp, err := client.SearchMovie(params)
	if err != nil {
		return nil, fmt.Errorf("搜索电影「%s」失败: %v", name, err)
	}
	movies := resp.Result
	// 排序：按年份降序
	sort.Slice(movies, func(i, j int) bool {
		// 年份降序
		// 取 release_date 格式 "YYYY-MM-DD"
		getDate := func(r themoviedb.SearchMovieResponse) string {
			if r.ReleaseDate != "" {
				return r.ReleaseDate
			}
			return "0000-00-00"
		}
		return strings.Compare(getDate(movies[i]), getDate(movies[j])) > 0
	})
	for _, movie := range movies {
		if utils.FuzzyMatching(name, movie.Title, movie.OriginalTitle) {
			info, err := GetMovieDetail(movie.ID)
			if err != nil {
				return nil, fmt.Errorf("获取电影「%d」详情失败: %v", movie.ID, err)
			}
			return info, nil
		}
	}
	return nil, fmt.Errorf("未找到电影「%s」的匹配结果", name)
}
