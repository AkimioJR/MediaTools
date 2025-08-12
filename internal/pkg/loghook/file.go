package loghook

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const chanSize = 500 // 缓冲通道大小

type FileLogsHook struct {
	logDir    string
	file      *os.File // 当前日志文件句柄
	day       int      // 用于记录当前日志文件的日期
	ch        chan *logrus.Entry
	formatter logrus.Formatter
	wg        sync.WaitGroup
}

func NewFileLogsHook(logDir string) (*FileLogsHook, error) {
	fh := FileLogsHook{
		logDir: logDir,
		ch:     make(chan *logrus.Entry, chanSize),
		formatter: &logrus.JSONFormatter{
			TimestampFormat:   time.DateTime, // 使用标准时间格式
			PrettyPrint:       true,          // 设置为 true 以启用格式化
			DisableHTMLEscape: true,          // 禁用 HTML 转义
		},
	}
	err := os.MkdirAll(logDir, 0755) // 确保日志目录存在
	if err != nil {
		return nil, fmt.Errorf("create log directory failed: %v", err)
	}
	go fh.writeLog()
	return &fh, nil
}

// Close 方法用于关闭文件和通道，确保资源被正确释放，避免协程泄漏
func (f *FileLogsHook) Close() {
	close(f.ch) // 关闭通道
	f.wg.Wait() // 等待协程退出
	if f.file != nil {
		f.file.Close()
		f.file = nil // 清空文件句柄
	}
}

func (f *FileLogsHook) ChangeLogDir(logDir string) error {
	err := os.MkdirAll(logDir, 0755) // 确保新日志目录存在
	if err != nil {
		return fmt.Errorf("create log directory failed: %v", err)
	}
	f.Close()
	f.logDir = logDir
	f.day = 0
	f.ch = make(chan *logrus.Entry, chanSize)
	go f.writeLog()
	return nil
}

// writeLog 是一个协程，用于异步写入日志到文件，避免并发写入问题
func (f *FileLogsHook) writeLog() {
	f.wg.Add(1)
	defer f.wg.Done()
	for entry := range f.ch {
		if f.day != entry.Time.Day() || f.file == nil {
			if f.file != nil {
				f.file.Close()
			}
			f.day = entry.Time.Day()
			var err error
			path := filepath.Join(f.logDir, entry.Time.Local().Format("2006-01-02")+".log")
			f.file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr, "open log file '%s' failed: %v", path, err)
				continue
			}
		}

		entryCopy := *entry    // 创建一个副本以避免修改原始条目
		entryCopy.Buffer = nil // 清空 Buffer 以确保使用 JSON 格式化器重新格式化

		logData, err := f.formatter.Format(&entryCopy)
		if err != nil {
			fmt.Fprintf(os.Stderr, "format log entry failed: %v", err)
			continue
		}

		if _, err := f.file.Write(logData); err != nil {
			fmt.Fprintf(os.Stderr, "write log to file failed: %v", err)
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
