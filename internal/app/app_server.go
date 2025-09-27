//go:build onlyServer

package app

var (
	Run                = runServer
	SupportDesktopMode = false
)
