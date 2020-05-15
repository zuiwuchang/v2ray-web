package web

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/version"
	"v2ray.com/core"
)

type _apiAPP struct {
}

func getSession(request *http.Request) (session *cookie.Session, e error) {
	c, e := request.Cookie(cookie.CookieName)
	if e != nil {
		if e == http.ErrNoCookie {
			e = nil
		}
		return
	}
	session, e = cookie.FromCookie(c.Value)
	return
}

// GetPost .
func GetPost(router *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) {
	router.GET(relativePath, handlers...)
	router.POST(relativePath, handlers...)
}

func (a *_apiAPP) RegisterTo(router *gin.RouterGroup) {
	r := router.Group(`/api/app`)
	GetPost(r, `/version`)
	GetPost(r, `/restore`)
	GetPost(r, `/login`)
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
		c.Error(e)
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
	e := c.BindJSON(&params)
	if e != nil {
		return
	}
	var mUser manipulator.User
	session, e := mUser.Login(params.Name, params.Password)
	if e != nil {
		return
	} else if session == nil {
		e = errors.New("name or password not match")
		return
	}
	val, e := session.Cookie()
	if e != nil {
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
