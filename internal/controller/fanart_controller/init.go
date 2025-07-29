package fanart_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/outbound"
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

	var opts []fanart.Options
	if config.Fanart.ApiURL != "" {
		opts = append(opts, fanart.CustomAPIURL(config.Fanart.ApiURL))
	}
	opts = append(opts, fanart.CustomHTTPClient(outbound.GetHTTPClient()))

	c, err := fanart.NewClient(config.Fanart.ApiKey, opts...)
	if err != nil {
		return err
	}
	client = c
	return nil
}
