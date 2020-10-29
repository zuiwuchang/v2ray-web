package v1

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/web"
)

// BaseURL .
const BaseURL = `/v1`

// Helper 一些其它的 api
type Helper struct {
	web.Helper
}

// Register impl IController
func (h Helper) Register(router *gin.RouterGroup) {
	r := router.Group(BaseURL)

	ms := []web.IHelper{
		// Debug{},
		Session{},
		Other{},
		// Users{},
		// Shells{},
		// FS{},
	}
	for _, m := range ms {
		m.Register(r)
	}
}
