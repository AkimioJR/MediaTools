package logging

import (
	"MediaTools/internal/config"

	"github.com/sirupsen/logrus"
)

func init() {
	setting := &serviceLoggerSetting{}
	logrus.SetFormatter(setting)
	logrus.AddHook(setting)
	logrus.SetReportCaller(true) // 启用调用者信息
	err := SetLogLevel(config.Log.Level)
	if err != nil {
		logrus.Errorf("设置日志级别失败，已设置为「%s」: %v", logrus.InfoLevel, err)
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func SetLogLevel(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)
	return nil
}
