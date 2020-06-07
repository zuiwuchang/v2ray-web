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
	"v2ray.com/ext/tools/conf/serial"
)

type _apiV2ray struct {
}

func (a *_apiV2ray) Init(router *gin.RouterGroup) {
	r := router.Group(`/v2ray`)
	GetPost(r, `/settings/get`, a.settingsGet)
	GetPost(r, `/settings/put`, a.settingsPut)
	GetPost(r, `/settings/test`, a.settingsTest)
	GetPost(r, `/subscription/list`, a.subscriptionList)
	GetPost(r, `/subscription/put`, a.subscriptionPut)
	GetPost(r, `/subscription/add`, a.subscriptionAdd)
	GetPost(r, `/subscription/remove`, a.subscriptionRemove)
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
	cnf, e := serial.LoadJSONConfig(&buffer)
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
func (a *_apiV2ray) subscriptionList(c *gin.Context) {
	var mSubscription manipulator.Subscription
	result, e := mSubscription.List()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (a *_apiV2ray) subscriptionPut(c *gin.Context) {
	var params data.Subscription
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	var mSubscription manipulator.Subscription
	e = mSubscription.Put(&params)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
}
func (a *_apiV2ray) subscriptionAdd(c *gin.Context) {
	var params data.Subscription
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSubscription manipulator.Subscription
	e = mSubscription.Add(&params)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, params.ID)
}
func (a *_apiV2ray) subscriptionRemove(c *gin.Context) {
	var params struct {
		ID uint64
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSubscription manipulator.Subscription
	e = mSubscription.Remove(params.ID)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, params.ID)
}
