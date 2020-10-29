package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/web"
	v1 "gitlab.com/king011/v2ray-web/web/api/v1"
)

// BaseURL request base url
const BaseURL = `/api`

// Helper path of /app
type Helper struct {
	web.Helper
}

// Register impl IHelper
func (h Helper) Register(router *gin.RouterGroup) {
	r := router.Group(BaseURL)

	ms := []web.IHelper{
		v1.Helper{},
	}
	for _, m := range ms {
		m.Register(r)
	}
}
