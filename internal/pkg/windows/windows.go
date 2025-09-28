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
	view   webview.WebView
	mutex  sync.RWMutex
	ch     chan bool // true -> 打开窗口 false -> 结束运行
}

func (w *Windows) createView() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.view == nil {
		w.view = webview.New(false)
		w.view.SetTitle(w.title)
		w.view.SetSize(w.width, w.height, webview.HintNone)
		w.view.Navigate(w.url)
		// fmt.Println("创建 view 成功")
	}

}

func (w *Windows) destroyView() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.view != nil {
		w.view.Destroy()
		w.view = nil
	}
}

func NewWindows(title, url string, width, height int) *Windows {
	w := Windows{
		title:  title,
		url:    url,
		width:  width,
		height: height,

		ch: make(chan bool),
	}
	return &w
}

func (w *Windows) Show() {
	if w.view == nil {
		w.ch <- true
		// fmt.Println("发送更新成功")
	}

}

func (w *Windows) Hide() {
	if w.view != nil {
		// fmt.Println("中断 view...")
		w.view.Terminate()
	}
}

func (w *Windows) Quit() {
	w.Hide()
	w.ch <- false
}

func (w *Windows) Run() {
	defer close(w.ch)
	defer w.destroyView()

	for {
		w.createView()
		w.view.Run()
		w.destroyView()

		// fmt.Println("view 运行停止，等待更新...")
		c := <-w.ch
		// fmt.Println("接受到：", c)
		if !c {
			// fmt.Println("退出")
			break
		}
	}
}
