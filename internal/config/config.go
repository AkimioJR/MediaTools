package config

import "github.com/spf13/viper"

var (
	Log  LogConfig
	TMDB TMDBConfig
)

func Init() error {
	var c Configuration
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&c); err != nil {
		return err
	}
	Log = c.Log
	TMDB = c.TMDB
	return nil
}
