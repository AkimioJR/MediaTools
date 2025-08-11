package loghook

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type FileLogsHook struct {
	LogDir string
}

func NewFileLogsHook(logDir string) *FileLogsHook {
	return &FileLogsHook{
		LogDir: logDir,
	}
}

func (f *FileLogsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *FileLogsHook) Fire(entry *logrus.Entry) error {
	if err := os.MkdirAll(f.LogDir, os.ModePerm); err != nil {
		return err
	}
	logFilePath := filepath.Join(f.LogDir, time.Now().Local().Format("2006-01-02")+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()

	formatter := &logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	}

	logData, err := formatter.Format(entry)
	if err != nil {
		return err
	}

	logFile.Write(logData)
	return nil
}

var _ logrus.Hook = (*FileLogsHook)(nil)
