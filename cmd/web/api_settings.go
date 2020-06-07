package web

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
)

type _apiSettings struct {
}

func (a *_apiSettings) Init(router *gin.RouterGroup) {
	r := router.Group(`/settings`)
	GetPost(r, `/get`, a.get)
	GetPost(r, `/put`, a.put)
}
func (a *_apiSettings) get(c *gin.Context) {
	var mSettings manipulator.Settings
	result, e := mSettings.Get()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, result)
	return
}
func (a *_apiSettings) put(c *gin.Context) {
	var params data.Settings
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.Put(&params)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	return
}
