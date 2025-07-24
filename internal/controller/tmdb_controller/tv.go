package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/themoviedb/v3"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// 根据电视剧的名称和季的年份及序号匹配TMDB
// name: 电视剧名称
// seasonYear: 季的年份
// seasonNumber: 季的序号
func SearchTVBySeason(name string, seasonYear int, seasonNumber int) (*schemas.MediaInfo, error) {
	logrus.Infof("正在搜索「%s (%d) 第 %d 季」...", name, seasonYear, seasonNumber)

	matchSeasonFN := func(tv themoviedb.SearchTVResponse) bool {
		detail, err := client.GetTVSerieDetail(tv.ID, nil)
		if err != nil {
			logrus.Warningf("获取电视剧「%d」详情失败: %v", tv.ID, err)
			return false
		}
		for _, season := range detail.Seasons {
			if season.SeasonNumber == seasonNumber && season.AirDate[:len("2025")] == strconv.Itoa(seasonYear) {
				return true
			}
		}
		return false
	}

	var params themoviedb.SearchTVSParams
	params.Query = name

	resp, err := client.SearchTV(params)
	if err != nil {
		return nil, fmt.Errorf("搜索电视剧「%s」失败: %v", name, err)
	}
	tvSeries := resp.Result
	if len(tvSeries) == 0 {
		return nil, fmt.Errorf("未找到电视剧「%s 第 %d 季 (%d)」的匹配结果", name, seasonNumber, seasonYear)
	}
	// 排序：按年份降序
	sort.Slice(tvSeries, func(i, j int) bool {
		// 年份降序
		// 取 release_date 格式 "YYYY-MM-DD"
		getDate := func(r themoviedb.SearchTVResponse) string {
			if r.FirstAirDate != "" {
				return r.FirstAirDate
			}
			return "0000-00-00"
		}
		return strings.Compare(getDate(tvSeries[i]), getDate(tvSeries[j])) > 0
	})

	for _, series := range tvSeries {
		seriesYear := series.FirstAirDate[:len("2025")]
		if utils.FuzzyMatching(name, series.Name, series.OriginalName) &&
			seriesYear == strconv.Itoa(seasonYear) {
			return GetTVSeriesDetail(series.ID)
		}
		names, err := getNames(series.ID, meta.MediaTypeTV)
		if err != nil {
			logrus.Warningf("获取电视剧「%d」的标题/译名失败: %v", series.ID, err)
			continue
		}

		if !utils.FuzzyMatching(name, names...) {
			continue
		}
		if matchSeasonFN(series) {
			info, err := GetTVSeriesDetail(series.ID)
			if err != nil {
				return nil, fmt.Errorf("获取电视剧「%d」详情失败: %v", series.ID, err)
			}
			return info, nil
		} else {
			logrus.Infof("电视剧「%s」(ID: %d) 不匹配指定季: %d-%d", series.Name, series.ID, seasonYear, seasonNumber)
		}
	}
	return nil, fmt.Errorf("未找到电视剧「%s」(季: %d-%d) 的匹配结果", name, seasonYear, seasonNumber)
}

func SearchTVByName(name string, year *int) (*schemas.MediaInfo, error) {
	var params themoviedb.SearchTVSParams
	params.Query = name
	if year != nil && *year > 0 {
		yearu32 := uint32(*year)
		params.Year = &yearu32
		logrus.Infof("正在搜索电视剧「%s (%d)」...", name, *year)
	} else {
		logrus.Infof("正在搜索电视剧「%s」...", name)
	}

	resp, err := client.SearchTV(params)
	if err != nil {
		return nil, fmt.Errorf("搜索电视剧「%s」失败: %v", name, err)
	}
	tvSeries := resp.Result
	// 排序：按年份降序
	sort.Slice(tvSeries, func(i, j int) bool {
		// 年份降序
		// 取 release_date 格式 "YYYY-MM-DD"
		getDate := func(r themoviedb.SearchTVResponse) string {
			if r.FirstAirDate != "" {
				return r.FirstAirDate
			}
			return "0000-00-00"
		}
		return strings.Compare(getDate(tvSeries[i]), getDate(tvSeries[j])) > 0
	})
	for _, series := range tvSeries {
		if utils.FuzzyMatching(name, series.Name, series.OriginalName) {
			info, err := GetTVSeriesDetail(series.ID)
			if err != nil {
				return nil, fmt.Errorf("获取电视剧「%d」详情失败: %v", series.ID, err)
			}
			return info, nil
		}
	}
	return nil, fmt.Errorf("未找到电视剧「%s」的匹配结果", name)
}
