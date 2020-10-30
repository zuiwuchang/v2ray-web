package web

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"v2ray.com/core"
)

type _apiV2ray struct {
}

func (a *_apiV2ray) Init(router *gin.RouterGroup) {
	r := router.Group(`/v2ray`)
	GetPost(r, `/settings/get`, a.settingsGet)
	GetPost(r, `/settings/put`, a.settingsPut)
	GetPost(r, `/settings/test`, a.settingsTest)
}

func (a *_apiV2ray) settingsGet(c *gin.Context) {
	var mSettings manipulator.Settings
	text, e := mSettings.GetV2ray()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, text)
}
func (a *_apiV2ray) settingsPut(c *gin.Context) {
	var params struct {
		Text string
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.PutV2ray(params.Text)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
}
func (a *_apiV2ray) settingsTest(c *gin.Context) {
	var params struct {
		Text string
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	outbound := &data.Outbound{
		Name:     "測試",
		Add:      "127.0.0.1",
		Port:     "1989",
		Net:      "tcp",
		Security: "auto",
		UserID:   "83b81e69-b1c7-077f-10d9-75b015b24651",
	}
	t := template.New("v2ray")
	t, e = t.Parse(params.Text)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	ctx, e := outbound.ToContext()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var buffer bytes.Buffer
	e = t.Execute(&buffer, ctx)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	// v2ray
	cnf, e := core.LoadConfig(`json`, `test.json`, &buffer)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	server, e := core.New(cnf)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	server.Close()
}
