package web

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/web"
	"gitlab.com/king011/v2ray-web/web/api"
	"gitlab.com/king011/v2ray-web/web/static"
	"gitlab.com/king011/v2ray-web/web/view"
)

func newGIN() (router *gin.Engine) {
	router = gin.Default()
	rs := []web.IHelper{
		view.Helper{},
		static.Helper{},
		api.Helper{},
	}
	for _, r := range rs {
		r.Register(&router.RouterGroup)
	}
	return
}
