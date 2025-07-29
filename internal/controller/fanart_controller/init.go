package fanart_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/fanart/v3"
	"sync"
)

var (
	client *fanart.FanartClient
	lock   sync.RWMutex
)

func Init() error {
	lock.Lock()
	defer lock.Unlock()

	c, err := fanart.NewClient(config.Fanart.ApiKey)
	if err != nil {
		return err
	}
	client = c
	return nil
}
