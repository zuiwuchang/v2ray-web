package web

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/speed"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// IModule .
type IModule interface {
	Init(router *gin.RouterGroup)
}

func getURL() (url string) {
	var mSettings manipulator.Settings
	result, e := mSettings.Get()
	if e == nil {
		url = strings.TrimSpace(result.URL)
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = speed.DefaultURL
		}
	} else {
		if ce := logger.Logger.Check(zap.WarnLevel, "get settings error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		url = speed.DefaultURL
	}
	return
}

// GetPost .
func GetPost(router *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) {
	router.GET(relativePath, handlers...)
	router.POST(relativePath, handlers...)
}

var pub = map[string]bool{
	`/`:                    true,
	`/index.html`:          true,
	`/angular`:             true,
	`/api/ws/proxy/status`: true,
	`/api/app/version`:     true,
	`/api/app/restore`:     true,
	`/api/app/login`:       true,
}

func checkRequest(c *gin.Context) {
	path := c.Request.URL.Path
	if strings.HasPrefix(path, `/angular`) ||
		pub[path] {
		return
	}
	session, e := getSession(c.Request)
	if e != nil {
		c.AbortWithError(http.StatusUnauthorized, e)
		return
	}
	if !session.Root {
		c.AbortWithError(http.StatusUnauthorized, e)
		return
	}
	return
}
func getSession(request *http.Request) (session *cookie.Session, e error) {
	c, e := request.Cookie(cookie.CookieName)
	if e != nil {
		if e == http.ErrNoCookie {
			e = nil
		}
		return
	}
	val, e := url.QueryUnescape(c.Value)
	if e != nil {
		return
	}
	session, e = cookie.FromCookie(val)
	return
}
