package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/web"
)

// Subscriptions 訂閱設定
type Subscriptions struct {
	web.Helper
}

// Register impl IHelper
func (h Subscriptions) Register(router *gin.RouterGroup) {
	r := router.Group(`subscriptions`)
	r.Use(h.CheckSession)

	r.GET(``, h.list)
	r.PUT(``, h.put)
	r.POST(``, h.add)
	r.DELETE(``, h.remove)
}
func (h Subscriptions) list(c *gin.Context) {
	var mSubscription manipulator.Subscription
	result, e := mSubscription.List()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, result)
}
func (h Subscriptions) put(c *gin.Context) {
	var obj data.Subscription
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}

	var mSubscription manipulator.Subscription
	e = mSubscription.Put(&obj)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Subscriptions) add(c *gin.Context) {
	var obj data.Subscription
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mSubscription manipulator.Subscription
	e = mSubscription.Add(&obj)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusCreated, obj.ID)
}
func (h Subscriptions) remove(c *gin.Context) {
	var obj struct {
		ID uint64 `form:"id" json:"id" xml:"id" yaml:"id" binding:"required"`
	}
	e := c.Bind(&obj)
	if e != nil {
		return
	}
	var mSubscription manipulator.Subscription
	e = mSubscription.Remove(obj.ID)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
