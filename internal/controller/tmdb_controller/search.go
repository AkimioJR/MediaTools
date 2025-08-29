package tmdb_controller

import (
	"MediaTools/internal/pkg/themoviedb/v3"
	"fmt"
	"iter"

	"github.com/sirupsen/logrus"
)

func SearchMovie(searchName string) (iter.Seq2[*themoviedb.SearchMovieResponse, error], error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("正在搜索电影「%s」...", searchName)
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

	return func(yield func(*themoviedb.SearchMovieResponse, error) bool) {
		for _, result := range resp.Result {
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
	}, nil
}

func SearchTV(searchName string) (iter.Seq2[*themoviedb.SearchTVResponse, error], error) {
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

	return func(yield func(*themoviedb.SearchTVResponse, error) bool) {
		for _, result := range resp.Result {
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
	}, nil
}

func SearchMulti(searchName string) (iter.Seq2[*themoviedb.SearchMultiResponse, error], error) {
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

	return func(yield func(*themoviedb.SearchMultiResponse, error) bool) {
		for _, res := range resp.Result {
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
	}, nil
}
