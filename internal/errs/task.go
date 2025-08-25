package errs

import "errors"

var (
	ErrTaskNotFound    = errors.New("task is not found")
	ErrTaskQueueClosed = errors.New("task queue is closed") // 任务队列已关闭
)
