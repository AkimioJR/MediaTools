package info

import "time"

var RuntimeAppStatus = RuntimeAppStatusInfo{
	DesktopMode: false,
	IsDev:       false,
	Port:        0,
	BootTime:    time.Now(),
}
