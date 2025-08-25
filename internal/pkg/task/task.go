package task

import (
	"context"
)

type TaskFunc func(ctx context.Context)

type Task struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	State TaskState `json:"state"`

	fn TaskFunc

	ctx    context.Context
	cancel context.CancelFunc
}
