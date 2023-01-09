package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xtls/xray-core/core"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/internal/net"
	"gitlab.com/king011/v2ray-web/template"
	"gitlab.com/king011/v2ray-web/web"
)

// V2ray 設定
type V2ray struct {
	web.Helper
}

// Register impl IHelper
func (h V2ray) Register(router *gin.RouterGroup) {
	r := router.Group(`v2ray`)
	r.Use(h.CheckSession)

	r.GET(``, h.get)
	r.PUT(``, h.put)
	r.POST(`test`, h.test)
	r.POST(`preview`, h.preview)
	r.GET(`default`, h.def)
}
func (h V2ray) get(c *gin.Context) {
	var mSettings manipulator.Settings
	text, e := mSettings.GetV2ray()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, text)
}
func (h V2ray) put(c *gin.Context) {
	var obj struct {
		Text string `form:"text" json:"text" xml:"text" yaml:"text" binding:"required"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.PutV2ray(obj.Text)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h V2ray) test(c *gin.Context) {
	var obj struct {
		Text string `form:"text" json:"text" xml:"text" yaml:"text" binding:"required"`
		URL  string `form:"url" json:"url" xml:"url" yaml:"url" binding:"required"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	protocol, result := net.AnalyzeString(obj.URL)
	if result == nil {
		h.NegotiateErrorString(c, http.StatusBadRequest, `not support proxy url`)
		return
	}
	outbound := &data.Outbound{
		Name:     result.Name,
		Add:      result.Add,
		Port:     result.Port,
		Host:     result.Host,
		TLS:      result.TLS,
		Net:      result.Net,
		Path:     result.Path,
		UserID:   result.UserID,
		AlterID:  result.AlterID,
		Security: result.Security,
		Level:    result.Level,
		Protocol: protocol,
		Flow:     result.Flow,
	}
	text, e := outbound.Render(obj.Text)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	// v2ray
	cnf, e := core.LoadConfig(`json`, strings.NewReader(text))
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	server, e := core.New(cnf)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	server.Close()
}
func (h V2ray) preview(c *gin.Context) {
	var obj struct {
		Text string `form:"text" json:"text" xml:"text" yaml:"text" binding:"required"`
		URL  string `form:"url" json:"url" xml:"url" yaml:"url" binding:"required"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	protocol, result := net.AnalyzeString(obj.URL)
	if result == nil {
		h.NegotiateErrorString(c, http.StatusBadRequest, `not support proxy url`)
		return
	}
	outbound := &data.Outbound{
		Name:     result.Name,
		Add:      result.Add,
		Port:     result.Port,
		Host:     result.Host,
		TLS:      result.TLS,
		Net:      result.Net,
		Path:     result.Path,
		UserID:   result.UserID,
		AlterID:  result.AlterID,
		Security: result.Security,
		Level:    result.Level,
		Protocol: protocol,
		Flow:     result.Flow,
	}

	text, e := outbound.Render(obj.Text)
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
func (h V2ray) def(c *gin.Context) {
	h.NegotiateData(c, http.StatusOK, gin.H{
		"text": template.Default,
	})
}
