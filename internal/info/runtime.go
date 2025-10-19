package info

import "time"

var RuntimeAppStatus = RuntimeAppStatusInfo{
	PlatformInfo: p,
	DesktopMode:  false,
	IsDev:        false,
	Port:         0,
	BootTime:     time.Now(),
}
