package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// 解析配置文件内容
func parseConfig(file *os.File) error {
	var config Configuration
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return fmt.Errorf("config parse error: %w", err)
	}

	config.applyConfig()
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
	err := os.MkdirAll(filepath.Dir(ConfigFile), 0755)
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
