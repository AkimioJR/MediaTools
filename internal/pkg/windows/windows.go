package windows

import (
	"sync"

	webview "github.com/webview/webview_go"
)

type Windows struct {
	title  string
	url    string
	width  int
	height int
	mutex  sync.RWMutex
	view   webview.WebView
	ch     chan bool // true -> 打开窗口 false -> 结束运行
}

func NewWindows(title, url string, width, height int) *Windows {
	w := Windows{
		title:  title,
		url:    url,
		width:  width,
		height: height,
		ch:     make(chan bool),
	}
	return &w
}

func (w *Windows) Close() {
	close(w.ch)
}

func (w *Windows) IsHide() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.view == nil
}

func (w *Windows) HideWindows() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.view != nil {
		w.view.Terminate()
	}
}

func (w *Windows) ShowWindow() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.view == nil {
		w.ch <- true
	}
}

func (w *Windows) Quit() {
	w.HideWindows()
	w.ch <- false
}

func (w *Windows) runWebView() {
	w.view = webview.New(false)
	w.view.SetTitle(w.title)
	w.view.SetSize(w.width, w.height, webview.HintNone)
	w.view.Navigate(w.url)
	defer func() {
		w.mutex.Lock()
		defer w.mutex.Unlock()
		w.view.Destroy()
		w.view = nil
	}()

	w.view.Run()
}

func (w *Windows) Run(fn func()) {

	for {
		w.runWebView()

		if fn != nil {
			fn()
		}

		// fmt.Println("view 运行停止，等待更新...")
		c := <-w.ch
		// fmt.Println("接受到：", c)
		if !c {
			// fmt.Println("退出")
			break
		}
	}
}
