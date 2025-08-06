package config

import (
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFile = "config/config.yaml"
)

var (
	Log        LogConfig
	TMDB       TMDBConfig
	Fanart     FanartConfig
	Media      MediaConfig
	CustomWord CustomWordConfig
)

func Init() error {
	var c Configuration
	viper.SetConfigFile(ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&c); err != nil {
		return err
	}

	Log = c.Log
	TMDB = c.TMDB
	Fanart = c.Fanart
	Media = c.Media
	CustomWord = c.CustomWord
	return nil
}

func WriteConfig() error {
	var c Configuration
	Log = c.Log
	TMDB = c.TMDB
	Fanart = c.Fanart
	Media = c.Media
	CustomWord = c.CustomWord
	file, err := os.OpenFile(ConfigFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := yaml.NewEncoder(file).Encode(&c); err != nil {
		return err
	}
	return nil
}
