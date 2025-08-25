package task

type TaskState uint8

const (
	TaskStatePending   TaskState = iota // 等待中
	TaskStateRunning                    // 运行中
	TaskStateCanceling                  // 取消中
)

func (ts TaskState) String() string {
	switch ts {
	case TaskStatePending:
		return "Pending"
	case TaskStateRunning:
		return "Running"
	case TaskStateCanceling:
		return "Canceling"
	default:
		return "UnknownTaskState"
	}
}

func (ts TaskState) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ts.String() + `"`), nil
}
