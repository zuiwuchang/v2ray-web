package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/speed"
)

type _apiProxy struct {
}

func (a *_apiProxy) Init(router *gin.RouterGroup) {
	r := router.Group(`/proxy`)
	GetPost(r, `/list`, a.list)
	GetPost(r, `/update`, a.update)
	GetPost(r, `/add`, a.add)
	GetPost(r, `/put`, a.put)
	GetPost(r, `/remove`, a.remove)
	GetPost(r, `/clear`, a.clear)
	GetPost(r, `/start`, a.start)
	GetPost(r, `/stop`, a.stop)
	GetPost(r, `/test`, a.test)
}
func (a *_apiProxy) list(c *gin.Context) {
	var mElement manipulator.Element
	element, subscription, e := mElement.List()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	c.JSON(http.StatusOK, &struct {
		Element      []*data.Element      `json:"element,omitempty"`
		Subscription []*data.Subscription `json:"subscription,omitempty"`
	}{
		element,
		subscription,
	})
	return
}
func (a *_apiProxy) update(c *gin.Context) {
	var params struct {
		ID uint64
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSubscription manipulator.Subscription
	info, e := mSubscription.Get(params.ID)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	outbounds, e := requestSubscription(info.URL)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	count := len(outbounds)
	if count == 0 {
		e = errors.New("outbound empty")
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mElement manipulator.Element
	result, e := mElement.Puts(info.ID, outbounds)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, result)
	return
}
func (a *_apiProxy) add(c *gin.Context) {
	var params struct {
		Subscription uint64
		Outbound     data.Outbound
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mElement manipulator.Element
	result, e := mElement.Add(params.Subscription, &params.Outbound)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, result)
	return
}
func (a *_apiProxy) put(c *gin.Context) {
	var params struct {
		ID           uint64
		Subscription uint64
		Outbound     data.Outbound
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mElement manipulator.Element
	e = mElement.Put(params.Subscription, params.ID, &params.Outbound)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	return
}

func (a *_apiProxy) remove(c *gin.Context) {
	var params struct {
		ID           uint64
		Subscription uint64
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mElement manipulator.Element
	e = mElement.Remove(params.Subscription, params.ID)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	return
}
func (a *_apiProxy) clear(c *gin.Context) {
	var params struct {
		Subscription uint64
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mElement manipulator.Element
	e = mElement.Clear(params.Subscription)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	return
}
func (a *_apiProxy) start(c *gin.Context) {
	var params data.Element
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	e = srv.Start(&params)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSettings manipulator.Settings
	e0 := mSettings.PutLast(&params)
	if e0 != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "save last status error"); ce != nil {
			ce.Write(
				zap.Error(e0),
			)
		}
	}
	return
}
func (a *_apiProxy) stop(c *gin.Context) {
	srv.Stop()
	return
}
func (a *_apiProxy) test(c *gin.Context) {
	var params data.Outbound
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	duration, e := speed.TestOne(&params, getURL())
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, duration.Milliseconds())
	return
}
