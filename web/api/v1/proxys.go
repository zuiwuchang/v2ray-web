package v1

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/internal/srv"
	"gitlab.com/king011/v2ray-web/logger"
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

	r.GET(``, h.list)
	r.DELETE(``, h.remove)
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
		ID           uint64 `form:"id" json:"id" xml:"id" yaml:"id" binding:"required"`
		Subscription uint64 `form:"subscription" json:"subscription" xml:"subscription" yaml:"subscription" binding:"required"`
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
