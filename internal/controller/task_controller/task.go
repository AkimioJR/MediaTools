package task_controller

import (
	"MediaTools/internal/pkg/task"
	"context"
)

var (
	c                 = context.Background() // 全局上下文
	transferTaskQueue = task.NewTaskQueue(c) // 转移任务队列
)
