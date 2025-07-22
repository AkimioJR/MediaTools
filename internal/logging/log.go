package logging

import (
	"github.com/sirupsen/logrus"
)

func init() {
	setting := &serviceLoggerSetting{}
	logrus.SetFormatter(setting)
	logrus.AddHook(setting)
	logrus.SetLevel(logrus.InfoLevel) // 默认日志级别为 Info
	logrus.SetReportCaller(true)      // 启用调用者信息
}
