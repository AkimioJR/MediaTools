package logging

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type RecentLogsHook struct {
	logs  []string
	size  uint
	index uint
	mu    sync.Mutex
}

func NewRecentLogsHook(size uint) *RecentLogsHook {
	return &RecentLogsHook{
		logs: make([]string, size),
		size: size,
	}
}

func (h *RecentLogsHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.logs[h.index] = entry.Message
	h.index = (h.index + 1) % h.size
	return nil
}

func (h *RecentLogsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 获取最新日志
func (h *RecentLogsHook) GetRecentLogs() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]string, 0, h.size)
	for i := range h.logs {
		idx := (h.index - 1 - uint(i) + h.size) % h.size
		if h.logs[idx] != "" {
			result = append(result, h.logs[idx])
		}
	}
	return result
}
