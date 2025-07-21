package limiter

import (
	"sync/atomic"
	"time"
)

type Limiter struct {
	d        time.Duration
	maxCount uint64
	count    atomic.Uint64
}

func NewLimiter(d time.Duration, maxCount uint64) *Limiter {
	l := &Limiter{
		d:        d,
		maxCount: maxCount,
	}
	go l.resetLoop()
	return l
}

func (l *Limiter) resetLoop() {
	ticker := time.NewTicker(l.d)
	defer ticker.Stop()
	for range ticker.C {
		l.count.Store(0)
	}
}

// Acquire 原子方式实现，超过限额时自旋等待
func (l *Limiter) Acquire() {
	for {
		cur := l.count.Load()
		if cur < l.maxCount {
			if l.count.CompareAndSwap(cur, cur+1) {
				return
			}
		} else {
			time.Sleep(10 * time.Millisecond) // 忙等，避免CPU空转
		}
	}
}

// NewLimitFunc 返回一个带限流功能的函数
func NewLimitFunc[In any, Out any](d time.Duration, maxCount uint64, fn func(In) Out) func(In) Out {
	l := NewLimiter(d, maxCount)
	return func(arg In) Out {
		l.Acquire()
		return fn(arg)
	}
}
