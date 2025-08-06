package library_controller

import "sync"

var (
	lock = sync.RWMutex{} // 锁定转移操作，防止并发冲突
)

func Init() error {
	return nil
}
