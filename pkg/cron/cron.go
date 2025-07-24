package cron

import "time"

type Cron struct {
	d        time.Duration
	stopChan chan struct{}
}

func NewCron(d time.Duration) *Cron {
	return &Cron{
		d:        d,
		stopChan: make(chan struct{}),
	}
}

func (c *Cron) Do(f func()) {
	go func() {
		ticker := time.NewTicker(c.d)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			case <-c.stopChan:
				return
			}
		}
	}()
}

func (c *Cron) Stop() {
	select {
	case _, _ = <-c.stopChan:
	default:
		c.stopChan <- struct{}{}
		close(c.stopChan)
	}
}
