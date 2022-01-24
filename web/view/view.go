package view

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/static"
	"gitlab.com/king011/v2ray-web/web"
)

// BaseURL request base url
const BaseURL = `view`

// Helper path of /app
type Helper struct {
	web.Helper
}

var zhHant = static.ZhHant()
var zhHans = static.ZhHans()

// Register impl IHelper
func (h Helper) Register(router *gin.RouterGroup) {
	router.GET(``, h.redirect)
	router.HEAD(``, h.redirect)
	router.GET(`index`, h.redirect)
	router.HEAD(`index`, h.redirect)
	router.GET(`index.html`, h.redirect)
	router.HEAD(`index.html`, h.redirect)
	router.GET(`view`, h.redirect)
	router.HEAD(`view`, h.redirect)
	router.GET(`view/`, h.redirect)
	router.HEAD(`view/`, h.redirect)

	r := router.Group(BaseURL)
	r.Use(h.Compression())

	r.GET(`:locale`, h.viewOrRedirect)
	r.HEAD(`:locale`, h.viewOrRedirect)
	r.GET(`:locale/*path`, h.view)
	r.HEAD(`:locale/*path`, h.view)
}
func (h Helper) redirect(c *gin.Context) {
	request := c.Request
	str := strings.ToLower(strings.TrimSpace(request.Header.Get(`Accept-Language`)))
	strs := strings.Split(str, `;`)
	str = strings.TrimSpace(strs[0])
	strs = strings.Split(str, `,`)
	str = strings.TrimSpace(strs[0])
	if strings.HasPrefix(str, `zh-`) {
		if strings.Contains(str, `cn`) || strings.Contains(str, `hans`) {
			c.Redirect(http.StatusFound, `/view/zh-Hans/`)
			return
		}
		c.Redirect(http.StatusFound, `/view/zh-Hant/`)
		return
	}
	c.Redirect(http.StatusFound, `/view/zh-Hant/`)
}
func (h Helper) viewOrRedirect(c *gin.Context) {
	var obj struct {
		Locale string `uri:"locale"`
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	if obj.Locale == "zh-Hant" {
		c.Redirect(http.StatusFound, `/view/zh-Hant/`)
	} else if obj.Locale == "zh-Hans" {
		c.Redirect(http.StatusFound, `/view/zh-Hans/`)
	} else {
		h.redirect(c)
	}
}
func (h Helper) view(c *gin.Context) {
	var obj struct {
		Locale string `uri:"locale" binding:"required"`
		Path   string `uri:"path" `
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	if obj.Locale == "zh-Hant" {
		c.Header("Cache-Control", "max-age=2419200")
		h.NegotiateFilesystem(c, zhHant, obj.Path)
	} else if obj.Locale == "zh-Hans" {
		c.Header("Cache-Control", "max-age=2419200")
		h.NegotiateFilesystem(c, zhHans, obj.Path)
	} else {
		h.NegotiateErrorString(c, http.StatusNotFound, `not support locale`)
	}
}
