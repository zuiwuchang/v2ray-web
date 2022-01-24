package web

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/static"
	"gitlab.com/king011/v2ray-web/web"
	"gitlab.com/king011/v2ray-web/web/api"
	"gitlab.com/king011/v2ray-web/web/view"
)

func newGIN() (router *gin.Engine) {
	router = gin.Default()
	// static
	router.Group(`static`).Use(web.Compression()).StaticFS(``, static.Static())
	// favicon.ico
	router.GET(`favicon.ico`, static.Favicon)
	router.HEAD(`favicon.ico`, static.Favicon)

	rs := []web.IHelper{
		view.Helper{},
		api.Helper{},
	}
	for _, r := range rs {
		r.Register(&router.RouterGroup)
	}
	return
}
