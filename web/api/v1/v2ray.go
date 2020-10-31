package v1

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/web"
	"v2ray.com/core"
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
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	outbound := &data.Outbound{
		Name:     "測試",
		Add:      "127.0.0.1",
		Port:     "1989",
		Net:      "tcp",
		Security: "auto",
		UserID:   "83b81e69-b1c7-077f-10d9-75b015b24651",
		Protocol: "vmess",
	}
	t := template.New("v2ray")
	t, e = t.Parse(obj.Text)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	ctx, e := outbound.ToContext()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	var buffer bytes.Buffer
	e = t.Execute(&buffer, ctx)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	str := buffer.String()
	// v2ray
	cnf, e := core.LoadConfig(`json`, `test.json`, &buffer)
	if e != nil {
		h.NegotiateErrorString(c, http.StatusInternalServerError, e.Error()+"\n"+str)
		return
	}
	server, e := core.New(cnf)
	if e != nil {
		h.NegotiateErrorString(c, http.StatusInternalServerError, e.Error()+"\n"+str)
		return
	}
	server.Close()
}
