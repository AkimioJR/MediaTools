package config

import (
	"fmt"
	"os"
	pathlib "path"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// 解析配置文件内容
func parseConfig(file *os.File) error {
	var c Configuration
	if err := yaml.NewDecoder(file).Decode(&c); err != nil {
		return fmt.Errorf("config parse error: %w", err)
	}
	c.check()
	c.applyConfig()
	return nil
}

// 初始化默认配置
func initDefaultConfig() error {
	defaultConfig.applyConfig()
	if err := WriteConfig(); err != nil {
		return fmt.Errorf("failed to create default config: %w", err)
	}
	return nil
}

// 应用配置到全局变量
func (c *Configuration) applyConfig() {
	Log = c.Log
	TMDB = c.TMDB
	Fanart = c.Fanart
	Storages = c.Storages
	Media = c.Media
}

func (c *Configuration) writeConfig() error {
	err := os.MkdirAll(pathlib.Dir(ConfigFile), 0755)
	if err != nil {
		return fmt.Errorf("create config directory error: %w", err)
	}
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

// 检查配置的完整性
func (c *Configuration) check() {
	if c.Log.ConsoleLevel == "" {
		logrus.Warning("日志终端输出级别未设置，使用默认配置")
		c.Log.ConsoleLevel = defaultConfig.Log.ConsoleLevel
	}

	if c.Log.FileLevel == "" {
		logrus.Warning("日志文件级别未设置，使用默认配置")
		c.Log.FileLevel = defaultConfig.Log.FileLevel
	}

	if c.Log.FileDir == "" {
		logrus.Warning("日志文件目录未设置，使用默认配置")
		c.Log.FileDir = defaultConfig.Log.FileDir
	}

	if c.TMDB.ApiURL == "" {
		logrus.Warning("TMDB API URL 配置未设置，使用默认配置")
		c.TMDB.ApiURL = defaultConfig.TMDB.ApiURL
	}

	if c.TMDB.ImageURL == "" {
		logrus.Warning("TMDB 图片 API URL 配置未设置，使用默认配置")
		c.TMDB.ImageURL = defaultConfig.TMDB.ImageURL
	}

	if c.Fanart.ApiURL == "" {
		logrus.Warning("Fanart API URL 配置未设置，使用默认配置")
		c.Fanart.ApiURL = defaultConfig.Fanart.ApiURL
	}

	if len(c.Storages) == 0 {
		logrus.Warning("存储配置未设置，使用默认配置")
		c.Storages = defaultConfig.Storages
	}

	if c.Media.Format == (FormatConfig{}) {
		logrus.Warning("媒体格式配置未设置，使用默认配置")
		c.Media.Format = defaultConfig.Media.Format
	}
}
