package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/web"
)

// 策略設定
type Strategy struct {
	web.Helper
}

// Register impl IHelper
func (h Strategy) Register(router *gin.RouterGroup) {
	r := router.Group(`strategys`)
	r.Use(h.CheckSession)

	r.GET(``, h.list)
	r.PUT(``, h.put)
}

func (h Strategy) list(c *gin.Context) {
	var mStrategy manipulator.Strategy
	result, e := mStrategy.List()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, result)
}
func (h Strategy) put(c *gin.Context) {
	var obj data.Strategy
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	obj.Name = strings.TrimSpace(obj.Name)
	if obj.Name == `` {
		h.NegotiateErrorString(c, http.StatusBadRequest, `name not supported: `+obj.Name)
		return
	}

	var mStrategy manipulator.Strategy
	e = mStrategy.Put(&obj)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
}
