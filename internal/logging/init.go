package logging

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/loghook"

	"github.com/sirupsen/logrus"
)

var (
	recentLogsHook = loghook.NewRecentLogsHook(100)
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

func SetLevel(level loghook.LogLevel) {
	logrus.SetLevel(level.ToLogrusLevel())
}

func GetRecentLogs() []loghook.LogDetail {
	return recentLogsHook.GetRecentLogs()
}
