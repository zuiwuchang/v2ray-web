package web

import (
	"io"
	"net/http"
	"sync/atomic"

	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/speed"
	"gitlab.com/king011/v2ray-web/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"v2ray.com/core/external/github.com/gorilla/websocket"
)

type _apiWebsocket struct {
	upgrader websocket.Upgrader
}

func (a *_apiWebsocket) Init(router *gin.RouterGroup) {
	a.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	r := router.Group(`/api/ws`)
	GetPost(r, `/app/logs`, a.appLogs)
	GetPost(r, `/proxy/test`, a.proxyTest)
	GetPost(r, `/proxy/status`, a.proxyStatus)
}
func (a *_apiWebsocket) appLogs(c *gin.Context) {
	ws, e := a.upgrader.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	defer ws.Close()
	cancel := make(chan struct{})
	var closed int32
	ch := make(chan string, 5)
	id := logs.AddListener(func(str string) {
		select {
		case ch <- str:
		default:
		}
	})
	go func() {
		var e error
		for {
			_, _, e = ws.ReadMessage()
			if e != nil {
				break
			}
		}
		if atomic.CompareAndSwapInt32(&closed, 0, 1) {
			close(cancel)
		}
	}()
	running := true
	for running {
		select {
		case <-cancel:
			running = false
		case str := <-ch:
			e := ws.WriteMessage(websocket.TextMessage, utils.StringToBytes(str))
			if e != nil {
				if atomic.CompareAndSwapInt32(&closed, 0, 1) {
					close(cancel)
				}
				running = false
			}
		}
	}
	logs.RemoveListener(id)
	return
}
func (a *_apiWebsocket) proxyStatus(c *gin.Context) {
	ws, e := a.upgrader.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	defer ws.Close()

	cancel := make(chan struct{})
	var closed int32
	ch := make(chan *ListenerStatus, 5)
	id := srv.AddListener(func(status *ListenerStatus) {
		select {
		case ch <- status:
		case <-cancel:
		}
	})
	go func() {
		var e error
		for {
			_, _, e = ws.ReadMessage()
			if e != nil {
				break
			}
		}
		if atomic.CompareAndSwapInt32(&closed, 0, 1) {
			close(cancel)
		}
	}()
	running := true
	for running {
		select {
		case <-cancel:
			running = false
		case status := <-ch:
			b, e := json.Marshal(status)
			if e == nil {
				e = ws.WriteMessage(websocket.TextMessage, b)
				if e != nil {
					if atomic.CompareAndSwapInt32(&closed, 0, 1) {
						close(cancel)
					}
					running = false
				}
			} else {
				if ce := logger.Logger.Check(zap.WarnLevel, "status marshal error"); ce != nil {
					ce.Write(
						zap.Error(e),
					)
				}
			}
		}
	}
	srv.RemoveListener(id)
}
func (a *_apiWebsocket) proxyTest(c *gin.Context) {
	ws, e := a.upgrader.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	defer ws.Close()

	ctx := speed.New(getURL())
	go ctx.Run()
	go func() {
		var e error
		var b []byte
		var t int
		for {
			t, b, e = ws.ReadMessage()
			if e != nil {
				ctx.Close()
				break
			} else if t != websocket.TextMessage {
				continue
			}
			msg := utils.BytesToString(b)
			if msg == "close" {
				ctx.CloseSend()
				break
			}
			var element speed.Element
			e = json.Unmarshal(b, &element)
			if e == nil {
				if !ctx.Send(&element) {
					break
				}
			} else {
				if ce := logger.Logger.Check(zap.WarnLevel, "unmarshal test element error"); ce != nil {
					ce.Write(
						zap.Error(e),
						zap.String("msg", msg),
					)
				}
			}
		}
	}()
	for {
		result := ctx.Get()
		if result == nil {
			break
		}
		b, e := json.Marshal(result)
		if e != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "marshal test result error"); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
		}
		e = ws.WriteMessage(websocket.TextMessage, b)
		if e != nil {
			ctx.Close()
			if e != io.EOF {
				if ce := logger.Logger.Check(zap.WarnLevel, "websocket closed"); ce != nil {
					ce.Write(
						zap.Error(e),
					)
				}
			}
			break
		}
	}
	e = ws.WriteMessage(websocket.TextMessage, []byte(`close`))
}
