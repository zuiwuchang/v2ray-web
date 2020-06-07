package web

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/version"
	"v2ray.com/core"
)

type _apiAPP struct {
}

func (a *_apiAPP) Init(router *gin.RouterGroup) {
	r := router.Group(`/app`)
	GetPost(r, `/version`, a.version)
	GetPost(r, `/restore`, a.restore)
	GetPost(r, `/login`, a.login)
	GetPost(r, `/logout`, a.logout)
}
func (a *_apiAPP) version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"platform": fmt.Sprintf("%v %v %v", runtime.GOOS, runtime.GOARCH, runtime.Version()),
		"tag":      version.Tag,
		"commit":   version.Commit,
		"date":     version.Date,
		"v2ray":    core.Version(),
	})
}
func (a *_apiAPP) restore(c *gin.Context) {
	session, e := getSession(c.Request)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, session)
}
func (a *_apiAPP) login(c *gin.Context) {
	var params struct {
		Name     string
		Password string
		Remember bool
	}
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mUser manipulator.User
	session, e := mUser.Login(params.Name, params.Password)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	} else if session == nil {
		e = errors.New("name or password not match")
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	val, e := session.Cookie()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	maxage := 0
	if params.Remember {
		maxage = int(cookie.MaxAge())
	}
	c.SetCookie(cookie.CookieName, val, maxage, `/`, ``, false, true)

	c.JSON(http.StatusOK, &session)
	return
}
func (a *_apiAPP) logout(c *gin.Context) {
	c.SetCookie(cookie.CookieName, `expired`, -1, `/`, ``, false, true)
	return
}
