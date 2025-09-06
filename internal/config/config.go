package config

import (
	"fmt"
	"os"

	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	RootDir      = "data"
	ConfigFile   = RootDir + "/" + "config.yaml"
	SQLiteDBFile = RootDir + "/" + "data.db"
)

var (
	DB       DataBaseConfig
	Log      LogConfig
	TMDB     TMDBConfig
	Fanart   FanartConfig
	Storages []StorageConfig
	Media    MediaConfig
	Version  = VersionInfo{
		AppVersion: appVersion,
		CommitHash: commitHash,
		BuildTime:  parseBuildTime(buildTime),
		GoVersion:  runtime.Version(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}
)

func Init() error {
	file, err := os.OpenFile(ConfigFile, os.O_RDONLY, 0644)
	switch {
	case err == nil:
		defer file.Close()
		return parseConfig(file)

	case os.IsNotExist(err):
		logrus.Warning("配置文件不存在，使用默认配置")
		return initDefaultConfig()

	default:
		return fmt.Errorf("打开配置文件失败: %w", err)
	}
}

func WriteConfig() error {
	var c = Configuration{
		DB:       DB,
		Log:      Log,
		TMDB:     TMDB,
		Fanart:   Fanart,
		Storages: Storages,
		Media:    Media,
	}
	return c.writeConfig()
}
