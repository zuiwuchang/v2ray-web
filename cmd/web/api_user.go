package web

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/manipulator"
)

type _apiUser struct {
}

func (a *_apiUser) Init(router *gin.RouterGroup) {
	r := router.Group(`/user`)
	GetPost(r, `/list`, a.list)
	GetPost(r, `/add`, a.add)
	GetPost(r, `/remove`, a.remove)
	GetPost(r, `/password`, a.password)
}
func (a *_apiUser) list(c *gin.Context) {
	var mUser manipulator.User
	result, e := mUser.List()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
func (a *_apiUser) add(c *gin.Context) {
	var params struct {
		Name     string
		Password string
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mUser manipulator.User
	e = mUser.Add(params.Name, params.Password)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	return
}
func (a *_apiUser) remove(c *gin.Context) {
	var params struct {
		Name string
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mUser manipulator.User
	e = mUser.Remove(params.Name)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	return
}
func (a *_apiUser) password(c *gin.Context) {
	var params struct {
		Name     string
		Password string
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mUser manipulator.User
	e = mUser.Password(params.Name, params.Password)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	return
}
