package loghook

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type FileLogsHook struct {
	LogDir    string
	file      *os.File
	day       int                // 用于跟踪当前日志文件的日期
	path      string             // 用于跟踪当前日志文件路径
	ch        chan *logrus.Entry // 用于异步写入日志
	wg        sync.WaitGroup     // 用于等待写入协程结束
	formatter logrus.Formatter   // 日志格式化器
}

func NewFileLogsHook(logDir string) *FileLogsHook {
	fh := FileLogsHook{
		LogDir: logDir,
		ch:     make(chan *logrus.Entry, 500), // 缓冲通道，用于异步写入日志
		formatter: &logrus.JSONFormatter{
			TimestampFormat: time.DateTime,
			PrettyPrint:     true, // 缩进
		},
	}
	go fh.writeLog() // 启动日志写入协程
	return &fh
}

// Close 方法用于关闭文件和通道，确保资源被正确释放，避免协程泄漏
func (f *FileLogsHook) Close() {
	close(f.ch)
	f.wg.Wait() // 等待写入协程完成
	if f.file != nil {
		f.file.Close()
	}
}

// writeLog 是一个协程，用于异步写入日志到文件，避免并发写入问题
func (f *FileLogsHook) writeLog() {
	for entry := range f.ch {
		f.wg.Add(1) // 增加等待组计数器
		if f.day != entry.Time.Day() || f.file == nil {
			if f.file != nil {
				f.file.Close()
			}
			f.day = entry.Time.Day()
			f.path = filepath.Join(f.LogDir, entry.Time.Local().Format("2006-01-02")+".log")
			var err error
			f.file, err = os.OpenFile(f.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr, "open log file '%s' failed: %v", f.path, err)
				f.wg.Done() // 结束当前日志条目的等待
				continue
			}
			logData, err := f.formatter.Format(entry)
			if err != nil {
				fmt.Fprintf(os.Stderr, "format log entry failed: %v", err)
				f.wg.Done() // 结束当前日志条目的等待
				continue
			}
			f.file.Write(logData)
			f.wg.Done() // 结束当前日志条目的等待
		}
	}
}

func (f *FileLogsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *FileLogsHook) Fire(entry *logrus.Entry) error {
	f.ch <- entry // 将日志条目发送到通道
	return nil
}

var _ logrus.Hook = (*FileLogsHook)(nil)
