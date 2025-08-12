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
	ch        chan logrus.Entry
	formatter logrus.Formatter
	wg        sync.WaitGroup
}

func NewFileLogsHook(logDir string) (*FileLogsHook, error) {
	fh := FileLogsHook{
		logDir: logDir,
		ch:     make(chan logrus.Entry, chanSize),
		formatter: &logrus.JSONFormatter{
			TimestampFormat: time.DateTime, // 使用标准时间格式
			PrettyPrint:     true,          // 设置为 true 以启用格式化
		},
	}
	go fh.writeLog()
	return &fh, nil
}

// Close 方法用于关闭文件和通道，确保资源被正确释放，避免协程泄漏
func (f *FileLogsHook) Close() {
	close(f.ch) // 关闭通道
	f.wg.Wait() // 等待协程退出
}

func (f *FileLogsHook) SetLogDir(logDir string) {
	f.Close()
	f.logDir = logDir
	f.ch = make(chan logrus.Entry, chanSize)
	go f.writeLog()
}

// writeLog 是一个协程，用于异步写入日志到文件，避免并发写入问题
func (f *FileLogsHook) writeLog() {
	var (
		file *os.File
		err  error
		day  = 0
	)

	f.wg.Add(1)
	defer func() {
		f.wg.Done()
		if file != nil {
			file.Close() // 确保文件在退出时被关闭
		}
	}()

	for entry := range f.ch {
		if f.logDir == "" { // 如果日志目录未设置，则跳过写入
			continue
		}
		if day != entry.Time.Day() || file == nil {
			if file != nil {
				file.Close()
			}
			day = entry.Time.Day()
			err = os.MkdirAll(f.logDir, 0755) // 确保新日志目录存在
			if err != nil {
				fmt.Fprintf(os.Stderr, "create log directory failed: %v", err)
			}
			path := filepath.Join(f.logDir, entry.Time.Local().Format("2006-01-02")+".log")
			file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr, "open log file '%s' failed: %v", path, err)
				continue
			}
		}

		entry.Buffer = nil                     // 清空 Buffer 以确保使用 JSON 格式化器重新格式化
		entry.Data["line"] = entry.Caller.Line // 添加行号到日志数据中

		logData, err := f.formatter.Format(&entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "format log entry failed: %v", err)
			continue
		}

		_, err = file.Write(logData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write log to file failed: %v", err)
		}
	}
}

func (f *FileLogsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *FileLogsHook) Fire(entry *logrus.Entry) error {
	f.ch <- *entry // 将日志条目发送到通道
	return nil
}

var _ logrus.Hook = (*FileLogsHook)(nil)
