package config

import (
	"MediaTools/internal/info"
	"fmt"

	"os"
	"path/filepath"
	"sync"

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
	DB = c.DB
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

var dbOnce sync.Once

// 检查配置的完整性
func (c *Configuration) check() {
	needSave := false
	dbOnce.Do(func() {
		if c.DB.Type == "" || c.DB.DSN == "" {
			logrus.Warning("数据库类型或连接字符串未设置，使用默认 SQLite 配置")
			c.DB.Type = defaultConfig.DB.Type
			c.DB.DSN = defaultConfig.DB.DSN
			needSave = true
		}
	})

	if c.Log.ConsoleLevel == "" {
		logrus.Warning("日志终端输出级别未设置，使用默认配置")
		c.Log.ConsoleLevel = defaultConfig.Log.ConsoleLevel
		needSave = true
	}

	if c.Log.FileLevel == "" {
		logrus.Warning("日志文件级别未设置，使用默认配置")
		c.Log.FileLevel = defaultConfig.Log.FileLevel
		needSave = true
	}

	if c.Log.FileDir == "" {
		logrus.Warning("日志文件目录未设置，使用默认配置")
		c.Log.FileDir = defaultConfig.Log.FileDir
		needSave = true
	}

	if c.TMDB.ApiURL == "" {
		logrus.Warning("TMDB API URL 配置未设置，使用默认配置")
		c.TMDB.ApiURL = defaultConfig.TMDB.ApiURL
		needSave = true
	}

	if c.TMDB.ImageURL == "" {
		logrus.Warning("TMDB 图片 API URL 配置未设置，使用默认配置")
		c.TMDB.ImageURL = defaultConfig.TMDB.ImageURL
		needSave = true
	}

	if c.Fanart.ApiURL == "" {
		logrus.Warning("Fanart API URL 配置未设置，使用默认配置")
		c.Fanart.ApiURL = defaultConfig.Fanart.ApiURL
		needSave = true
	}

	if len(c.Storages) == 0 {
		logrus.Warning("存储配置未设置，使用默认配置")
		c.Storages = defaultConfig.Storages
		needSave = true
	}

	if c.Media.Format == (FormatConfig{}) {
		logrus.Warning("媒体格式配置未设置，使用默认配置")
		c.Media.Format = defaultConfig.Media.Format
		needSave = true
	}

	if needSave {
		logrus.Info("需要更新配置文件")
		if err := c.writeConfig(); err != nil {
			logrus.Errorf("保存配置文件失败: %v", err)
		} else {
			logrus.Info("配置文件已更新")
		}
	}
}

func getExecDir() string {
	execPath, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("无法获取可执行文件路径: %v", err))
	}
	return filepath.Dir(execPath)
}

func getDataPath() string {
	if info.Version.SupportDesktopMode {
		return filepath.Join(getExecDir(), "data")
	} else {
		return "data"
	}
}

func getLogsPath() string {
	if info.Version.SupportDesktopMode {
		return filepath.Join(getExecDir(), "logs")
	} else {
		return "logs"
	}
}
