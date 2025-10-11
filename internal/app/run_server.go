//go:build !desktop
// +build !desktop

package app

import "MediaTools/internal/info"

func init() {
	info.Version.SupportDesktopMode = false
}

func Run() {
	runServer()
}
