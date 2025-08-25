package task

import (
	"MediaTools/internal/errs"
	"context"
	"sync"

	"github.com/google/uuid"
)

type TaskQueue struct {
	taskMap  sync.Map // key: id string; value task *Task
	taskChan chan *Task

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTaskQueue(c context.Context) *TaskQueue {
	chanSize := 10
	workerNum := 5

	ctx, cancel := context.WithCancel(c)
	tq := TaskQueue{
		taskChan: make(chan *Task, chanSize),
		ctx:      ctx,
		cancel:   cancel,
	}
	tq.wg.Add(workerNum)
	for range workerNum {
		go tq.run()
	}

	return &tq
}

// 消费者函数
func (tq *TaskQueue) run() {
	defer tq.wg.Done()

	for {
		select {
		case task := <-tq.taskChan:
			if task.State == TaskStatePending { // 仅任务处于等待时才执行
				task.State = TaskStateRunning
				task.fn(task.ctx)
				tq.taskMap.Delete(task.ID)
			}

		case <-tq.ctx.Done():
			return
		}
	}
}

// 关闭任务队列
// 该操作会取消改队列下的所有任务
func (tq *TaskQueue) Close() {
	tq.taskMap.Clear()
	tq.cancel()
	tq.wg.Wait()
}

// 生产者函数
// 向任务队列中添加任务
func (tq *TaskQueue) SubmitTask(name string, fn TaskFunc) *Task {
	var (
		id string
		ok bool
	)

	for ok {
		id := uuid.New().String()
		_, ok = tq.taskMap.Load(id)
	}

	ctx, cancel := context.WithCancel(tq.ctx)
	task := &Task{
		ID:     id,
		Name:   name,
		fn:     fn,
		ctx:    ctx,
		cancel: cancel,
	}
	tq.taskMap.Store(id, task)
	tq.taskChan <- task
	return task
}

// 获取任务
func (tq *TaskQueue) GetTask(id string) (*Task, error) {
	value, ok := tq.taskMap.Load(id)
	if !ok {
		return nil, errs.ErrTaskNotFound
	}
	return value.(*Task), nil
}

// 取消任务
func (tq *TaskQueue) CancelTask(id string) (*Task, error) {
	task, err := tq.GetTask(id)
	if err != nil {
		return nil, err
	}

	task.State = TaskStateCanceling
	defer func() {
		task.cancel()
		tq.taskMap.Delete(id)
	}()
	return task, nil
}

func (tq *TaskQueue) IterTasks(yield func(task *Task) bool) {
	tq.taskMap.Range(func(key, value any) bool {
		task := value.(*Task)
		if !yield(task) {
			return false
		} else {
			return true
		}
	})
}
