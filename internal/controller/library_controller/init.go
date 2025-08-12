package library_controller

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	lock = sync.RWMutex{} // 锁定转移操作，防止并发冲突
)

func Init() error {
	lock.Lock()
	defer lock.Unlock()

	logrus.Info("开始初始化 Library Controller...")

	logrus.Info("Library Controller 初始化完成")
	return nil

}
