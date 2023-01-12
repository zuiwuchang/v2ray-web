package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/web"
)

// Settings 系統設定
type Settings struct {
	web.Helper
}

// Register impl IHelper
func (h Settings) Register(router *gin.RouterGroup) {
	r := router.Group(`settings`)
	r.Use(h.CheckSession)

	r.GET(``, h.get)
	r.PUT(``, h.put)
}
func (h Settings) get(c *gin.Context) {
	var obj struct {
		Strategy bool `form:"strategy"`
	}
	e := h.BindQuery(c, &obj)
	if e != nil {
		return
	}

	var mSettings manipulator.Settings
	result, e := mSettings.Get()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	if !obj.Strategy {
		h.NegotiateData(c, http.StatusOK, result)
		return
	}
	var mStrategy manipulator.Strategy
	strategys, e := mStrategy.List()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, map[string]any{
		`settings`:  result,
		`strategys`: strategys,
	})
}
func (h Settings) put(c *gin.Context) {
	var obj data.Settings
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.Put(&obj)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
