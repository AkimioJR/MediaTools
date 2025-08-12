package tmdb_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/outbound"
	"MediaTools/internal/pkg/themoviedb/v3"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	client *themoviedb.Client
	lock   = sync.RWMutex{}
)

func Init() error {
	lock.Lock()
	defer lock.Unlock()

	logrus.Info("开始初始化 TMDB Controller...")
	var opts []themoviedb.ClientOptions
	if config.TMDB.Language != "" {
		opts = append(opts, themoviedb.CustomLanguage(config.TMDB.Language))
	}
	if config.TMDB.IncludeImageLanguage != "" {
		opts = append(opts, themoviedb.CustomImageLanguage(config.TMDB.IncludeImageLanguage))
	}
	if config.TMDB.ApiURL != "" {
		opts = append(opts, themoviedb.CustomAPIURL(config.TMDB.ApiURL))
	}
	if config.TMDB.ImageURL != "" {
		opts = append(opts, themoviedb.CustomImageURL(config.TMDB.ImageURL))
	}
	opts = append(opts, themoviedb.CustomHTTPClient(outbound.GetHTTPClient()))

	c, err := themoviedb.NewClient(config.TMDB.ApiKey, opts...)
	if err != nil {
		return err
	}
	client = c
	logrus.Info("TMDB Controller 初始化完成")
	return nil

}
