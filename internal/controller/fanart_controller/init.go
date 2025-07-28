package fanart_controller

import (
	"MediaTools/internal/pkg/fanart/v3"
	"sync"
)

var (
	client *fanart.FanartClient
	lock   sync.RWMutex
)

func Init(apikey string) error {
	lock.Lock()
	defer lock.Unlock()

	c, err := fanart.NewClient(apikey)
	if err != nil {
		return err
	}
	client = c
	return nil
}
