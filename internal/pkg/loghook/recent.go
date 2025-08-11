package loghook

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type RecentLogsHook struct {
	logs  []LogDetail
	size  uint
	index uint
	lock  sync.RWMutex
}

func NewRecentLogsHook(size uint) *RecentLogsHook {
	return &RecentLogsHook{
		logs: make([]LogDetail, size),
		size: size,
	}
}

func (h *RecentLogsHook) Fire(entry *logrus.Entry) error {
	h.lock.Lock()
	defer h.lock.Unlock()
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
	h.lock.RLock()
	defer h.lock.RUnlock()

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
