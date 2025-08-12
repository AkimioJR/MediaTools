package logging

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/loghook"
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	historyLogsHook *loghook.MemoryHistoryHook
	fileHook        *loghook.FileLogsHook
)

func init() {
	var err error
	historyLogsHook = loghook.NewMemoryHistoryHook(100) // 初始化内存历史日志钩子
	fileHook, err = loghook.NewFileLogsHook("")         // 初始化文件日志钩子
	if err != nil {
		panic("初始化文件日志钩子失败: " + err.Error())
	}
	logrus.SetReportCaller(true) // 启用调用者信息

	f := &Formater{}
	logrus.SetFormatter(f)

	logrus.AddHook(fileHook)
	logrus.AddHook(historyLogsHook)
}

func Init() error {
	logrus.Debug("初始化日志系统...")
	fileHook.SetLogDir(config.Log.Path) // 设置日志目录
	err := SetLevel(config.Log.Level)   // 设置日志级别
	if err != nil {
		return fmt.Errorf("设置日志级别失败: %w", err)
	}
	logrus.Debug("日志系统初始化完成")
	return nil
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
