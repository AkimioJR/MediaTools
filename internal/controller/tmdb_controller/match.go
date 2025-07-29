package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/themoviedb/v3"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// 根据名称同时查询电影和电视剧，没有类型也没有年份时使用
// name 查询的名称
// mType 优先返回的媒体类型，如果为 nil 则电影优先
func MatchMulti(name string, mType *meta.MediaType) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	var page uint32 = 1
	var params themoviedb.SearchMultiParams
	params.Query = name
	params.Page = &page

	resp, err := client.SearchMulti(params)
	if err != nil {
		return nil, fmt.Errorf("综合查询 「%s」失败: %v", name, err)
	}
	results := make([]themoviedb.SearchMultiResponse, 0, len(resp.Result))
	results = append(results, resp.Result...)
	if resp.TotalPages > 1 {
		for i := 2; i <= int(resp.TotalPages); i++ {
			page = uint32(i)
			params.Page = &page
			pageResp, err := client.SearchMulti(params)
			if err != nil {
				return nil, fmt.Errorf("综合查询「%s」第 %d 页失败: %v", name, i, err)
			}
			results = append(results, pageResp.Result...)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("未找到综合查询结果「%s」", name)
	}

	// 排序：电影在前，按年份降序
	sort.Slice(results, func(i, j int) bool {
		if results[i].MediaType != results[j].MediaType {
			if mType != nil { // 如果指定了类型、对应类型优先
				switch *mType {
				case meta.MediaTypeMovie:
					if results[i].MediaType != results[j].MediaType {
						return results[i].MediaType == "movie"
					}
				case meta.MediaTypeTV:
					if results[i].MediaType != results[j].MediaType {
						return results[i].MediaType == "tv"
					}
				}
			} else { // 否则电影优先
				return results[i].MediaType == "movie"
			}
		}

		// 如果类型一致，年份新的优先
		// 取 release_date 格式 "YYYY-MM-DD"
		getDate := func(r themoviedb.SearchMultiResponse) string {
			if r.ReleaseDate != "" {
				return r.ReleaseDate
			}
			return "0000-00-00"
		}
		return strings.Compare(getDate(results[i]), getDate(results[j])) > 0
	})

	var info *schemas.MediaInfo

	for _, result := range results {
		var err error
		switch result.MediaType {
		case "movie":
			if utils.FuzzyMatching(name, result.Title, result.OriginalTitle) {
				info, err = GetMovieDetail(result.ID)
				if err != nil {
					logrus.Warningf("获取电影「%d」详情失败: %v", result.ID, err)
					continue
				}
				break
			}
			names, err := getNames(result.ID, meta.MediaTypeMovie)
			if err != nil {
				logrus.Errorf("获取电影「%d」的标题/译名失败: %v", result.ID, err)
				continue
			}
			if utils.FuzzyMatching(name, names...) {
				info, err = GetMovieDetail(result.ID)
				if err != nil {
					logrus.Warningf("获取电影「%d」详情失败: %v", result.ID, err)
					continue
				}
				goto match
			}
		case "tv":
			if utils.FuzzyMatching(name, result.Title, result.OriginalTitle) {
				info, err = GetTVSeriesDetail(result.ID)
				if err != nil {
					logrus.Warningf("获取电视剧「%d」详情失败: %v", result.ID, err)
					continue
				}
				break
			}
			names, err := getNames(result.ID, meta.MediaTypeTV)
			if err != nil {
				logrus.Errorf("获取电视剧「%d」的标题/译名失败: %v", result.ID, err)
				continue
			}
			if utils.FuzzyMatching(name, names...) {
				info, err = GetTVSeriesDetail(result.ID)
				if err != nil {
					logrus.Warningf("获取电视剧「%d」详情失败: %v", result.ID, err)
					continue
				}
				goto match
			}
		default:
			logrus.Warningf("不支持的媒体类型: %s", result.MediaType)
			continue
		}
	}
	if info == nil {
		return nil, fmt.Errorf("综合查询「%s」未找到结果", name)
	}
match:
	logrus.Infof("综合查询「%s」结果: %d (%s)", name, info.TMDBID, info.MediaType.String())
	return info, nil
}

// 搜索 TMDB 中的媒体信息，匹配返回一条尽可能正确的信息
// name 媒体名称
// mType 媒体类型
// year 年份，如要是季集需要是首播年份(可选)
// seasonYear 当前季集年份(可选)
// seasonNumber 当前季集数(可选)
func Match(name string, mType meta.MediaType, year *int, seasonYear *int, seasonNumber *int) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	switch mType {
	case meta.MediaTypeMovie:
		return SearchMovieByName(name, year)

	case meta.MediaTypeTV:
		if seasonYear != nil && seasonNumber != nil {
			return SearchTVBySeason(name, *seasonYear, *seasonNumber)
		} else {
			return SearchTVByName(name, year)
		}

	default:
		logrus.Warningf("未指定媒体类型，尝试综合查询「%s」", name)
		info, err := MatchMulti(name, nil)
		if err != nil {
			return nil, fmt.Errorf("综合查询「%s」失败: %v", name, err)
		}
		return GetInfo(info.TMDBID, &mType)
	}
}

// 搜索 TMDB 中的媒体信息，匹配返回一条尽可能正确的信息
// 如果匹配失败，尝试移除年份重新查询
// name 媒体名称
// mType 媒体类型
// year 年份，如要是季集需要是首播年份(可选)
// seasonYear 当前季集年份(可选)
// seasonNumber 当前季集数(可选)
func MatchWithFallback(name string, mType meta.MediaType, year *int, seasonYear *int, seasonNumber *int) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	// 1. 首次严格匹配
	info, err := Match(name, mType, year, seasonYear, seasonNumber)
	if err == nil {
		return info, nil
	}

	// 2. 如果有年份/季年份/季号，尝试去掉这些条件再查一次
	if utils.CheckValue(year) || utils.CheckValue(seasonYear) || utils.CheckValue(seasonNumber) {
		logrus.Warningf("匹配 %s「%s」失败: %v，尝试移除年份/季年份/季号重新查询", mType, name, err)
		info, err = Match(name, mType, nil, nil, nil)
		if err == nil {
			return info, nil
		}
	}

	// 3. 兜底：多类型模糊匹配
	logrus.Warningf("匹配 %s「%s」失败: %v，尝试综合模糊匹配", mType, name, err)
	info, err = MatchMulti(name, &mType)
	if err == nil {
		return info, nil
	}

	// 4. 移除掉类型再次匹配
	logrus.Warningf("综合模糊匹配匹配 %s「%s」失败: %v，移除类型再次尝试匹配", mType, name, err)
	info, err = MatchMulti(name, nil)
	if err == nil {
		return info, nil
	}

	// 5. 最终失败
	logrus.Errorf("综合匹配「%s」失败: %v", name, err)
	return nil, err
}
