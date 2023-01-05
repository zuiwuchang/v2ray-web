package speed

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitlab.com/king011/v2ray-web/template"

	"github.com/xtls/xray-core/core"
	"gitlab.com/king011/v2ray-web/db/data"
	"golang.org/x/net/proxy"
)

// DefaultURL .
const DefaultURL = `https://www.youtube.com/`

var mutex sync.Mutex

// TestOne .
func TestOne(outbound *data.Outbound, url string) (duration time.Duration, e error) {
	mutex.Lock()
	duration, e = testOne(outbound, 10000-1989, url)
	mutex.Unlock()
	return
}
func testOne(outbound *data.Outbound, port int, url string) (duration time.Duration, e error) {
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
	text, e := outbound.Render(template.Proxy)
	if e != nil {
		return
	}
	// v2ray
	cnf, e := core.LoadConfig(`json`, strings.NewReader(text))
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
	e = requestURL(port, url)
	if e != nil {
		return
	}
	duration = time.Since(last)
	return
}
func requestURL(port int, url string) (e error) {
	client := &http.Client{}
	var dialer proxy.Dialer
	// fmt.Println(port)
	dialer, e = proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%v", port), nil, proxy.Direct)
	if e != nil {
		return
	}
	client.Timeout = time.Second * 5
	client.Transport = &http.Transport{
		Dial: dialer.Dial,
	}
	response, e := client.Get(url)
	if e != nil {
		return
	}
	defer response.Body.Close()
	_, e = io.ReadAll(response.Body)
	return
}
