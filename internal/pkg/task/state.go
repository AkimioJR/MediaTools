package task

type TaskState uint8

const (
	TaskStatePending   TaskState = iota // 等待中
	TaskStateRunning                    // 运行中
	TaskStateCanceling                  // 取消中
)
