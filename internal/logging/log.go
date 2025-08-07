package logging

import (
	"MediaTools/internal/config"

	"github.com/sirupsen/logrus"
)

var (
	recentLogsHook = NewRecentLogsHook(10)
)

func Init() error {
	setting := &serviceLoggerSetting{}
	logrus.SetFormatter(setting)
	logrus.AddHook(setting)
	logrus.AddHook(recentLogsHook)
	logrus.SetReportCaller(true) // 启用调用者信息
	SetLevel(config.Log.Level)
	return nil
}

func SetLevel(level config.LogLevel) {
	logrus.SetLevel(level.ToLogrusLevel())
}

func GetRecentLogs() []string {
	return recentLogsHook.GetRecentLogs()
}
