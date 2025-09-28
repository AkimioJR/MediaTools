//go:build !desktop
// +build !desktop

package app

var (
	SupportDesktopMode = false
)

func Run() {
	runServer()
}
