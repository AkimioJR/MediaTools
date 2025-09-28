//go:build onlyServer

package app

var (
	SupportDesktopMode = false
)

func Run() {
	runServer()
}
