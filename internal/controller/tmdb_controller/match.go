package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/themoviedb/v3"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"context"
	"fmt"
	"iter"

	"github.com/sirupsen/logrus"
)

func MatchMovie(ctx context.Context, searchName string) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("正在搜索「%s」...", searchName)

	var firstResult *themoviedb.SearchMovieResponse
	var results iter.Seq2[*themoviedb.SearchMovieResponse, error] = func(yield func(*themoviedb.SearchMovieResponse, error) bool) {
		var page uint32 = 1
		var params = themoviedb.SearchMovieParams{
			Query: searchName,
			Page:  &page,
		}
		resp, err := client.SearchMovie(params)
		if err != nil {
			yield(nil, fmt.Errorf("搜索电影「%s」失败: %v", searchName, err))
			return
		}
		if resp.TotalResults == 0 {
			yield(nil, fmt.Errorf("未找到电影「%s」", searchName))
			return
		}
		for _, result := range resp.Result {
			if firstResult == nil {
				firstResult = &result // 保存第一条结果
			}
			if !yield(&result, nil) {
				return
			}
		}
		if resp.TotalPages > 1 {
			for page = 2; page <= uint32(resp.TotalPages); page++ {
				resp, err = client.SearchMovie(params)
				if err != nil {
					if !yield(nil, fmt.Errorf("搜索电影「%s」第 %d 页失败: %v", searchName, page, err)) {
						return
					}
				}
				for _, result := range resp.Result {
					if !yield(&result, nil) {
						return
					}
				}
			}
		}
	}

	for result, err := range results {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("搜索电影任务被取消: %w", ctx.Err())

		default:
		}
		if err != nil {
			logrus.Warning(err)
			continue // 如果搜索失败，尝试下一个结果
		}

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
	logrus.Warningf("未找到电影「%s」的匹配项，返回第一项: %s", searchName, firstResult.Title)
	return GetInfo(firstResult.ID, meta.MediaTypeMovie)
}

func MatchTV(ctx context.Context, searchName string) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("正在搜索「%s」...", searchName)
	var firstResult *themoviedb.SearchTVResponse
	var results iter.Seq2[*themoviedb.SearchTVResponse, error] = func(yield func(*themoviedb.SearchTVResponse, error) bool) {
		var page uint32 = 1
		var params = themoviedb.SearchTVSParams{
			Query: searchName,
			Page:  &page,
		}
		resp, err := client.SearchTV(params)
		if err != nil {
			yield(nil, fmt.Errorf("搜索电视剧「%s」失败: %v", searchName, err))
			return
		}
		if resp.TotalResults == 0 {
			yield(nil, fmt.Errorf("未找到电视剧「%s」", searchName))
			return
		}
		for _, result := range resp.Result {
			if firstResult == nil {
				firstResult = &result // 保存第一条结果
			}
			if !yield(&result, nil) {
				return
			}
		}
		if resp.TotalPages > 1 {
			for page = 2; page <= uint32(resp.TotalPages); page++ {
				resp, err = client.SearchTV(params)
				if err != nil {
					if !yield(nil, fmt.Errorf("搜索电视剧「%s」第 %d 页失败: %v", searchName, page, err)) {
						return
					}
				}
				for _, result := range resp.Result {
					if !yield(&result, nil) {
						return
					}
				}
			}
		}
	}

	for result, err := range results {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("搜索电视剧任务被取消: %w", ctx.Err())

		default:
		}
		if err != nil {
			logrus.Warning(err)
			continue // 如果搜索失败，尝试下一个结果
		}

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
	logrus.Warningf("未找到电视剧「%s」的匹配项，返回第一项: %s", searchName, firstResult.Name)
	return GetInfo(firstResult.ID, meta.MediaTypeTV)
}

func MatchMulti(ctx context.Context, searchName string) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("正在综合搜索「%s」...", searchName)
	var firstResult *themoviedb.SearchMultiResponse
	var results iter.Seq2[*themoviedb.SearchMultiResponse, error] = func(yield func(*themoviedb.SearchMultiResponse, error) bool) {
		var page uint32 = 1
		var params = themoviedb.SearchMultiParams{
			Query: searchName,
			Page:  &page,
		}
		resp, err := client.SearchMulti(params)
		if err != nil {
			yield(nil, fmt.Errorf("综合搜索「%s」失败: %v", searchName, err))
			return
		}
		if resp.TotalResults == 0 {
			yield(nil, fmt.Errorf("未找到综合搜索结果「%s」", searchName))
			return
		}
		for _, res := range resp.Result {
			if firstResult == nil {
				firstResult = &res // 保存第一条结果
			}
			if !yield(&res, nil) {
				return
			}
		}
		if resp.TotalPages > 1 {
			for page = 2; page <= uint32(resp.TotalPages); page++ {
				resp, err = client.SearchMulti(params)
				if err != nil {
					if !yield(nil, fmt.Errorf("综合搜索「%s」第 %d 页失败: %v", searchName, page, err)) {
						return
					}
				}
				for _, res := range resp.Result {
					if !yield(&res, nil) {
						return
					}
				}
			}
		}
	}

	for result, err := range results {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("综合搜索任务被取消: %w", ctx.Err())
		default:
		}

		if err != nil {
			logrus.Warning(err)
			continue // 如果搜索失败，尝试下一个结果
		}

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

	logrus.Warningf("未找到综合搜索「%s」的匹配项，返回第一项: %s", searchName, firstResult.Title)
	mediaType := parseType(firstResult.MediaType)
	if mediaType == meta.MediaTypeUnknown {
		logrus.Warningf("综合搜索结果「%s」的媒体类型(%s)未知", firstResult.Title, firstResult.MediaType)
	}
	return GetInfo(firstResult.ID, mediaType)
}
