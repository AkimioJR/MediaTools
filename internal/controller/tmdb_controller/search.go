package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/themoviedb/v3"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"fmt"

	"github.com/sirupsen/logrus"
)

func SearchMovie(searchName string) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("正在搜索「%s」...", searchName)
	var page uint32 = 1
	var params = themoviedb.SearchMovieParams{
		Query: searchName,
		Page:  &page,
	}

	resp, err := client.SearchMovie(params)
	if err != nil {
		return nil, fmt.Errorf("搜索电影「%s」失败: %v", searchName, err)
	}
	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("未找到电影「%s」", searchName)
	}

	var results = make([]themoviedb.SearchMovieResponse, 0, resp.TotalResults)
	results = append(results, resp.Result...)

	if resp.TotalPages > 1 {
		for page = 2; page <= uint32(resp.TotalPages); page++ {
			resp, err = client.SearchMovie(params)
			if err != nil {
				return nil, fmt.Errorf("搜索电影「%s」第 %d 页失败: %v", searchName, page, err)
			}
			results = append(results, resp.Result...)
		}
	}

	for _, result := range results {
		if utils.FuzzyMatching(searchName, result.Title, result.OriginalTitle) {
			logrus.Infof("匹配电影「%s」(TMDB ID: %d)", result.Title, result.ID)

			return GetInfo(result.ID, meta.MediaTypeMovie)
		}
		names, err := getNames(result.ID, meta.MediaTypeMovie)
		if err != nil {
			logrus.Warnf("获取电影「%s(%d)」的其他名称失败: %v", result.Title, result.ID, err)
			continue
		}
		if utils.FuzzyMatching(searchName, names...) {
			logrus.Infof("匹配电影「%s」(TMDB ID: %d) 别名", result.Title, result.ID)

			return GetInfo(result.ID, meta.MediaTypeMovie)
		}
	}
	logrus.Warningf("未找到电影「%s」的匹配项，返回第一项", searchName)
	return GetInfo(results[0].ID, meta.MediaTypeMovie)
}

func SearchTV(searchName string) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("正在搜索「%s」...", searchName)
	var page uint32 = 1
	var params = themoviedb.SearchTVSParams{
		Query: searchName,
		Page:  &page,
	}

	resp, err := client.SearchTV(params)
	if err != nil {
		return nil, fmt.Errorf("搜索电视剧「%s」失败: %v", searchName, err)
	}
	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("未找到电视剧「%s」", searchName)
	}

	var results = make([]themoviedb.SearchTVResponse, 0, resp.TotalResults)
	results = append(results, resp.Result...)

	if resp.TotalPages > 1 {
		for page = 2; page <= uint32(resp.TotalPages); page++ {
			resp, err = client.SearchTV(params)
			if err != nil {
				return nil, fmt.Errorf("搜索电视剧「%s」第 %d 页失败: %v", searchName, page, err)
			}
			results = append(results, resp.Result...)
		}
	}

	for _, result := range results {
		if utils.FuzzyMatching(searchName, result.Name, result.OriginalName) {
			logrus.Infof("匹配电视剧「%s」(TMDB ID: %d)", result.Name, result.ID)

			return GetInfo(result.ID, meta.MediaTypeTV)
		}
		names, err := getNames(result.ID, meta.MediaTypeTV)
		if err != nil {
			logrus.Warnf("获取电视剧「%s(%d)」的其他名称失败: %v", result.Name, result.ID, err)
			continue
		}
		if utils.FuzzyMatching(searchName, names...) {
			logrus.Infof("匹配电视剧「%s」(TMDB ID: %d) 别名", result.Name, result.ID)

			return GetInfo(result.ID, meta.MediaTypeTV)
		}
	}
	logrus.Warningf("未找到电视剧「%s」的匹配项，返回第一项", searchName)
	return GetInfo(results[0].ID, meta.MediaTypeTV)
}

func SearchMulti(searchName string) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("正在综合搜索「%s」...", searchName)
	var page uint32 = 1
	var params = themoviedb.SearchMultiParams{
		Query: searchName,
		Page:  &page,
	}
	resp, err := client.SearchMulti(params)
	if err != nil {
		return nil, fmt.Errorf("综合搜索「%s」失败: %v", searchName, err)
	}
	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("未找到综合搜索结果「%s」", searchName)
	}
	var results = make([]themoviedb.SearchMultiResponse, 0, resp.TotalResults)
	results = append(results, resp.Result...)
	if resp.TotalPages > 1 {
		for page = 2; page <= uint32(resp.TotalPages); page++ {
			resp, err = client.SearchMulti(params)
			if err != nil {
				return nil, fmt.Errorf("综合搜索「%s」第 %d 页失败: %v", searchName, page, err)
			}
			results = append(results, resp.Result...)
		}
	}
	for _, result := range results {
		mediaType := parseType(result.MediaType)
		if mediaType == meta.MediaTypeUnknown {
			logrus.Warningf("综合搜索结果「%s」的媒体类型(%s)未知，跳过", result.Title, result.MediaType)
			continue
		}

		if utils.FuzzyMatching(searchName, result.Title, result.OriginalTitle, result.Name, result.OriginalName) {
			logrus.Infof("匹配综合搜索结果「%s」(Type: %s TMDB ID: %d)", result.Title, mediaType, result.ID)
			return GetInfo(result.ID, mediaType)
		}
		names, err := getNames(result.ID, mediaType)
		if err != nil {
			logrus.Warnf("获取综合搜索结果「%s(Type: %s TMDB ID: %d)」的其他名称失败: %v", result.Title, mediaType, result.ID, err)
			continue
		}
		if utils.FuzzyMatching(searchName, names...) {
			logrus.Infof("匹配综合搜索结果「%s」(Type: %s TMDB ID: %d) 别名", result.Title, mediaType, result.ID)

			return GetInfo(result.ID, mediaType)
		}
	}
	result := results[0]
	mediaType := parseType(result.MediaType)
	if mediaType == meta.MediaTypeUnknown {
		logrus.Warningf("综合搜索结果「%s」的媒体类型(%s)未知", result.Title, result.MediaType)
	}
	return GetInfo(result.ID, mediaType)
}
