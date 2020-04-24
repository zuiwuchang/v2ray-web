package speed

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"golang.org/x/net/proxy"
	"v2ray.com/core"
	"v2ray.com/ext/tools/conf/serial"
)

// Context 速度測試 環境
type Context struct {
	cancel chan struct{}
	closed int32
	ch     chan *Element
	result chan *Result
}

// New .
func New() *Context {
	return &Context{
		cancel: make(chan struct{}),
		ch:     make(chan *Element),
		result: make(chan *Result),
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
	// 查詢可用 tcp 端口
	var target int
	for i := 0; i < 2000; i++ {
		str := fmt.Sprintf("127.0.0.1:%v", port+i)
		l, e := net.Listen("tcp", str)
		if e != nil {
			continue
		}
		l.Close()
		target = port + i
		break
	}
	if target == 0 {
		e = fmt.Errorf("not found idle port")
		return
	}
	ctx, e := element.Outbound.ToContext()
	if e != nil {
		return
	}
	t := template.New("v2ray")
	t, e = t.Parse(fmt.Sprintf(templateText, target))
	if e != nil {
		return
	}
	var buffer bytes.Buffer
	e = t.Execute(&buffer, ctx)
	if e != nil {
		return
	}
	// v2ray
	cnf, e := serial.LoadJSONConfig(&buffer)
	if e != nil {
		return
	}
	server, e := core.New(cnf)
	if e != nil {
		return
	}
	defer server.Close()
	e = server.Start()
	if e != nil {
		return
	}
	last := time.Now()
	e = s.http(element, target)
	if e != nil {
		return
	}
	result = &Result{
		ID:       element.ID,
		Status:   StatusOk,
		Duration: time.Now().Sub(last).Milliseconds(),
	}
	return
}
func (s *Context) http(element *Element, port int) (e error) {
	client := &http.Client{}
	var dialer proxy.Dialer
	dialer, e = proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%v", port), nil, proxy.Direct)
	if e != nil {
		return
	}
	client.Timeout = time.Second * 5
	client.Transport = &http.Transport{
		Dial: dialer.Dial,
	}
	response, e := client.Get("https://www.youtube.com/")
	if e != nil {
		return
	}
	if response.Body != nil {
		ioutil.ReadAll(response.Body)
	}
	return
}
