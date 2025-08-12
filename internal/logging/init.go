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
	fileHook.ChangeLogDir(config.Log.Path) // 设置日志目录
	return SetLevel(config.Log.Level)
}

func SetLevel(level string) error {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(logrusLevel)
	return nil
}

func GetRecentLogs() []loghook.LogDetail {
	return historyLogsHook.GetRecentLogs()
}
