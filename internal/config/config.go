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
	Log    LogConfig
	TMDB   TMDBConfig
	Fanart FanartConfig
	Media  MediaConfig
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
