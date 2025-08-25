package task_controller

import "MediaTools/internal/pkg/task"

func SubmitTransferTask(name string, fn task.TaskFunc) *task.Task {
	return transferTaskQueue.SubmitTask(name, fn)
}

func GetTransferTask(id string) (*task.Task, error) {
	return transferTaskQueue.GetTask(id)
}

func CancelTransferTask(id string) (*task.Task, error) {
	return transferTaskQueue.CancelTask(id)
}

func IterTransferTasks(yield func(t *task.Task) bool) {
	transferTaskQueue.IterTasks(yield)
}
