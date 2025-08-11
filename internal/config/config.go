package config

import (
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFile = "config/config.yaml"
)

var (
	Log     LogConfig
	TMDB    TMDBConfig
	Fanart  FanartConfig
	Media   MediaConfig
	Version = VersionInfo{
		AppVersion: appVersion,
		CommitHash: commitHash,
		BuildDate:  parseBuildTime(buildDate),
		GoVersion:  runtime.Version(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}
)

func Init() error {
	var c Configuration
	file, err := os.OpenFile(ConfigFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&c); err != nil {
		return err
	}

	Log = c.Log
	TMDB = c.TMDB
	Fanart = c.Fanart
	Media = c.Media
	return nil
}

func WriteConfig() error {
	var c Configuration
	c.Log = Log
	c.TMDB = TMDB
	c.Fanart = Fanart
	c.Media = Media
	file, err := os.Create(ConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := yaml.NewEncoder(file).Encode(&c); err != nil {
		return err
	}
	return nil
}
