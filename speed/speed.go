package speed

import (
	"sync"
	"sync/atomic"
)

// Context 速度測試 環境
type Context struct {
	cancel chan struct{}
	closed int32
	ch     chan *Element
	result chan *Result
	url    string
}

// New .
func New(url string) *Context {
	return &Context{
		cancel: make(chan struct{}),
		ch:     make(chan *Element),
		result: make(chan *Result),
		url:    url,
	}
}

// Close 停止測試 釋放資源
func (s *Context) Close() {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		close(s.cancel)
	}
}

// CloseSend 關閉 生產者
func (s *Context) CloseSend() {
	close(s.ch)
}

// Send 投遞一個 測試項目
func (s *Context) Send(element *Element) (ok bool) {
	select {
	case s.ch <- element:
		ok = true
	case <-s.cancel:
	}
	return
}

// Get 返回 響應數據
func (s *Context) Get() (result *Result) {
	select {
	case result = <-s.result:
	case <-s.cancel:
	}
	return
}

// Run 運行 測試
func (s *Context) Run() {
	var wait sync.WaitGroup
	wait.Add(10)
	for i := 0; i < 10; i++ {
		go s.run(&wait, 10000+1989+64+i*2000)
	}
	wait.Wait()
	s.Close()
}
func (s *Context) run(wait *sync.WaitGroup, port int) {
	for element := range s.ch {
		result := &Result{
			ID:     element.ID,
			Status: StatusRunning,
		}
		if !s.response(result) {
			break
		}

		result, e := s.do(element, port)
		if e != nil {
			result = &Result{
				ID:     element.ID,
				Status: StatusError,
				Error:  e.Error(),
			}
		}
		if result != nil {
			if !s.response(result) {
				break
			}
		}
	}
	wait.Done()

}
func (s *Context) response(result *Result) (ok bool) {
	select {
	case s.result <- result:
		ok = true
	case <-s.cancel:
	}
	return
}
func (s *Context) do(element *Element, port int) (result *Result, e error) {
	duration, _, e := testOne(&element.Outbound, port, s.url)
	if e != nil {
		return
	}
	result = &Result{
		ID:       element.ID,
		Status:   StatusOk,
		Duration: duration.Milliseconds(),
	}
	return
}
