package logging

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/loghook"

	"github.com/sirupsen/logrus"
)

var (
	historyLogsHook = loghook.NewMemoryHistoryHook(100)
	fileHook        = loghook.NewFileLogsHook("logs")
)

func init() {
	logrus.SetReportCaller(true) // 启用调用者信息

	f := &Formater{}
	logrus.SetFormatter(f)

	logrus.AddHook(fileHook)
	logrus.AddHook(historyLogsHook)
}
func Init() error {
	fileHook.LogDir = config.Log.Path
	SetLevel(config.Log.Level)
	return nil
}

func SetLevel(level loghook.LogLevel) {
	logrus.SetLevel(level.ToLogrusLevel())
}

func GetRecentLogs() []loghook.LogDetail {
	return historyLogsHook.GetRecentLogs()
}
