package transfer_controller

import "sync"

var (
	transferLock = sync.Mutex{} // 锁定转移操作，防止并发冲突
)

func Init() error {
	return nil
}
