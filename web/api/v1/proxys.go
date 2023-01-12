package v1

import (
	"io"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xtls/xray-core/core"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/internal/net"
	"gitlab.com/king011/v2ray-web/internal/srv"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/speed"
	"gitlab.com/king011/v2ray-web/utils"
	"gitlab.com/king011/v2ray-web/web"
	"go.uber.org/zap"
)

// Proxys 代理
type Proxys struct {
	web.Helper
}

// Register impl IHelper
func (h Proxys) Register(router *gin.RouterGroup) {
	r := router.Group(`proxys`)
	r.Use(h.CheckSession)

	r.GET(`status/websocket`, h.status)
	r.GET(`test/websocket`, h.test)

	r.GET(``, h.list)
	r.DELETE(``, h.remove)
	r.POST(``, h.add)
	r.PUT(``, h.put)
	r.POST(`start`, h.start)
	r.POST(`stop`, h.stop)
	r.DELETE(`clear`, h.clear)
	r.POST(`update`, h.update)
	r.POST(`test`, h.testOne)
	r.POST(`preview`, h.preview)
}

func (h Proxys) status(c *gin.Context) {
	ws, e := h.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	defer ws.Close()

	cancel := make(chan struct{})
	var closed int32
	ch := make(chan *srv.ListenerStatus, 5)
	id := srv.AddListener(func(status *srv.ListenerStatus) {
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
	json := h.JSON()
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
func (h Proxys) list(c *gin.Context) {
	var mElement manipulator.Element
	element, subscription, e := mElement.List()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}

	h.NegotiateData(c, http.StatusOK, &struct {
		Element      []*data.Element      `json:"element,omitempty" xml:"element,omitempty" yaml:"element,omitempty"`
		Subscription []*data.Subscription `json:"subscription,omitempty" xml:"subscription,omitempty" yaml:"subscription,omitempty"`
	}{
		element,
		subscription,
	})
}
func (h Proxys) remove(c *gin.Context) {
	var obj struct {
		ID           uint64 `form:"id" json:"id" xml:"id" yaml:"id"`
		Subscription uint64 `form:"subscription" json:"subscription" xml:"subscription" yaml:"subscription"`
	}
	e := c.Bind(&obj)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	e = mElement.Remove(obj.Subscription, obj.ID)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Proxys) start(c *gin.Context) {
	var obj struct {
		data.Element
		Strategy string `json:"strategy,omitempty" xml:"strategy,omitempty" yaml:"strategy,omitempty"`
	}
	e := c.Bind(&obj)
	if e != nil {
		return
	}

	text, e := srv.StartStrategy(&obj.Element, obj.Strategy)
	if e != nil {
		if text == `` {
			h.NegotiateError(c, http.StatusInternalServerError, e)
		} else {
			h.NegotiateData(c, http.StatusOK, gin.H{
				"text":  text,
				"error": e.Error(),
			})
		}
		return
	}
	var mSettings manipulator.Settings
	e0 := mSettings.PutLast(&obj.Element)
	if e0 != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "save last status error"); ce != nil {
			ce.Write(
				zap.Error(e0),
			)
		}
	}
	c.Status(http.StatusNoContent)
}
func (h Proxys) stop(c *gin.Context) {
	srv.Stop()
	c.Status(http.StatusNoContent)
}
func (h Proxys) clear(c *gin.Context) {
	var obj struct {
		Subscription uint64 `form:"subscription" json:"subscription" xml:"subscription" yaml:"subscription"`
	}
	e := c.Bind(&obj)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	e = mElement.Clear(obj.Subscription)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Proxys) add(c *gin.Context) {
	var obj struct {
		Subscription uint64        `form:"subscription" json:"subscription" xml:"subscription" yaml:"subscription"`
		Outbound     data.Outbound `form:"outbound" json:"outbound" xml:"outbound" yaml:"outbound"`
	}
	e := c.Bind(&obj)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	result, e := mElement.Add(obj.Subscription, &obj.Outbound)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusCreated, result)
}
func (h Proxys) put(c *gin.Context) {
	var obj struct {
		ID           uint64        `form:"id" json:"id" xml:"id" yaml:"id"`
		Subscription uint64        `form:"subscription" json:"subscription" xml:"subscription"`
		Outbound     data.Outbound `form:"outbound" json:"outbound" xml:"outbound" yaml:"outbound"`
	}
	e := c.Bind(&obj)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	e = mElement.Put(obj.Subscription, obj.ID, &obj.Outbound)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Proxys) update(c *gin.Context) {
	var obj struct {
		ID uint64 `form:"id" json:"id" xml:"id" yaml:"id"`
	}
	e := c.Bind(&obj)
	if e != nil {
		return
	}
	var mSubscription manipulator.Subscription
	info, e := mSubscription.Get(obj.ID)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	outbounds, e := net.RequestSubscription(info.URL)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	count := len(outbounds)
	if count == 0 {
		h.NegotiateErrorString(c, http.StatusInternalServerError, "outbound empty")
		return
	}
	var mElement manipulator.Element
	result, e := mElement.Puts(info.ID, outbounds)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, result)
}
func (h Proxys) testOne(c *gin.Context) {
	var obj data.Outbound
	e := c.Bind(&obj)
	if e != nil {
		return
	}
	duration, e := speed.TestOne(&obj, h.getURL())
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, duration.Milliseconds())
}
func (h Proxys) preview(c *gin.Context) {
	var obj struct {
		data.Outbound
		Strategy string `json:"strategy,omitempty" xml:"strategy,omitempty" yaml:"strategy,omitempty"`
	}
	e := c.Bind(&obj)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	text, e := mSettings.GetV2ray()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	var mStrategy manipulator.Strategy
	strategy, e := mStrategy.Value(obj.Strategy)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}

	text, e = obj.RenderStrategy(text, strategy)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}

	// v2ray
	cnf, e := core.LoadConfig(`json`, strings.NewReader(text))
	if e != nil {
		h.NegotiateData(c, http.StatusOK, gin.H{
			"text":  text,
			"error": e.Error(),
		})
		return
	}
	server, e := core.New(cnf)
	if e != nil {
		h.NegotiateData(c, http.StatusOK, gin.H{
			"text":  text,
			"error": e.Error(),
		})
		return
	}
	server.Close()
	h.NegotiateData(c, http.StatusOK, gin.H{
		"text": text,
	})
}
func (h Proxys) getURL() (url string) {
	var mSettings manipulator.Settings
	result, e := mSettings.Get()
	if e == nil {
		url = strings.TrimSpace(result.URL)
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = speed.DefaultURL
		}
	} else {
		if ce := logger.Logger.Check(zap.WarnLevel, "get settings error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		url = speed.DefaultURL
	}
	return
}
func (h Proxys) test(c *gin.Context) {
	ws, e := h.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	defer ws.Close()

	ctx := speed.New(h.getURL())
	json := h.JSON()
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
	ws.WriteMessage(websocket.TextMessage, []byte(`close`))
}
