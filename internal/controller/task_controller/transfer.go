package task_controller

import "MediaTools/internal/pkg/task"

func SubmitTransferTask(name string, fn task.TaskFunc) *task.Task {
	return transferTaskQueue.SubmitTask(name, fn)
}

func CancelTransferTask(id string) error {
	return transferTaskQueue.CancelTask(id)
}
