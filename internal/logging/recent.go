package logging

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type LogDetail struct {
	Level   logrus.Level `json:"level"`   // 日志级别
	Message string       `json:"message"` // 日志消息
	Time    time.Time    `json:"time"`    // 日志时间
	Caller  string       `json:"caller"`  // 日志调用者
}

type RecentLogsHook struct {
	logs  []LogDetail
	size  uint
	index uint
	mu    sync.Mutex
}

func NewRecentLogsHook(size uint) *RecentLogsHook {
	return &RecentLogsHook{
		logs: make([]LogDetail, size),
		size: size,
	}
}

func (h *RecentLogsHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.logs[h.index] = LogDetail{ // 直接存储结构体值
		Level:   entry.Level,
		Message: entry.Message,
		Time:    entry.Time,
		Caller:  entry.Caller.Function,
	}
	h.index = (h.index + 1) % h.size
	return nil
}

func (h *RecentLogsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 获取最新日志
func (h *RecentLogsHook) GetRecentLogs() []LogDetail {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]LogDetail, 0, h.size)
	for i := uint(0); i < h.size; i++ {
		idx := (h.index - 1 - i + h.size) % h.size
		if h.logs[idx].Time.IsZero() { // 跳过未初始化的条目（通过 Time 零值检测）
			continue
		}
		result = append(result, h.logs[idx])
	}
	return result
}
