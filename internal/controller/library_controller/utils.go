package library_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"slices"
	"strings"
)

func MatchLibrary(fi *storage.StorageFileInfo) *config.LibraryConfig {
	lock.RLock()
	defer lock.RUnlock()

	for i, lib := range config.Media.Libraries {
		if lib.SrcType == fi.StorageType &&
			strings.HasPrefix(fi.Path, lib.SrcPath) {
			return &config.Media.Libraries[i]
		}
	}
	return nil
}

func MatchCategory(cs []Category, countries []string, language string, genreIDs []int) string {
	lock.RLock()
	defer lock.RUnlock()

	var match bool

	for _, c := range cs {
		if c.OriginalCountries != nil {
			match = false
			for _, country := range countries {
				if slices.Contains(c.OriginalCountries, country) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		if c.OriginalLanguages != nil {
			match = false
			if slices.Contains(c.OriginalLanguages, language) {
				match = true
			}
			if !match {
				continue
			}
		}
		if c.GenreIDs != nil {
			match = false
			for _, genreID := range genreIDs {
				if slices.Contains(c.GenreIDs, genreID) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		return c.Name
	}
	return "未分类"
}

func GenFloder(libConfig *config.LibraryConfig, info *schemas.MediaInfo) []string {
	lock.RLock()
	defer lock.RUnlock()

	var floderNames []string
	if libConfig.OrganizeByType {
		switch info.MediaType {
		case meta.MediaTypeMovie:
			floderNames = append(floderNames, "电影")
		case meta.MediaTypeTV:
			floderNames = append(floderNames, "电视剧")
		}
	}
	if libConfig.OrganizeByCategory {
		var categoryStr string
		switch info.MediaType {
		case meta.MediaTypeMovie:
			var genres []int
			var countries []string

			for _, country := range info.TMDBInfo.MovieInfo.ProductionCountries {
				countries = append(countries, country.Iso31661)
			}
			language := info.TMDBInfo.MovieInfo.OriginalLanguage
			for _, genre := range info.TMDBInfo.MovieInfo.Genres {
				genres = append(genres, genre.ID)
			}
			categoryStr = MatchCategory(categoryConfig.MovieCategories, countries, language, genres)
		case meta.MediaTypeTV:
			var genres []int
			var countries []string

			for _, country := range info.TMDBInfo.TVInfo.SerieInfo.ProductionCountries {
				countries = append(countries, country.Iso31661)
			}
			language := info.TMDBInfo.TVInfo.SerieInfo.OriginalLanguage
			for _, genre := range info.TMDBInfo.TVInfo.SerieInfo.Genres {
				genres = append(genres, genre.ID)
			}
			categoryStr = MatchCategory(categoryConfig.TVCategories, countries, language, genres)
		}
		if categoryStr != "" {
			floderNames = append(floderNames, categoryStr)
		}
	}

	return floderNames
}
