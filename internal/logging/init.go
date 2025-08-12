package logging

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/loghook"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	historyLogsHook *loghook.MemoryHistoryHook
	fileHook        *loghook.FileLogsHook
)

func init() {
	var err error
	historyLogsHook = loghook.NewMemoryHistoryHook(100) // 初始化内存历史日志钩子
	formatter := logrus.JSONFormatter{
		TimestampFormat: time.DateTime, // 设置时间戳格式
		PrettyPrint:     true,          // 启用美化输出
	}
	fileHook, err = loghook.NewFileLogsHook("", loghook.WithFormatter(&formatter)) // 初始化文件日志钩子
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
	err := SetLevel(config.Log.ConsoleLevel) // 设置日志级别
	if err != nil {
		return fmt.Errorf("设置终端日志级别失败: %w", err)
	}
	err = SetFileLevel(config.Log.FileLevel) // 设置文件日志级别
	if err != nil {
		return fmt.Errorf("设置文件日志级别失败: %w", err)
	}
	SetLogDir(config.Log.FileDir) // 设置日志目录
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

func SetFileLevel(level string) error {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	fileHook.SetLevel(logrusLevel)
	return nil
}

func SetLogDir(dir string) {
	fileHook.SetLogDir(dir)
}

func GetRecentLogs() []loghook.LogDetail {
	return historyLogsHook.GetRecentLogs()
}
