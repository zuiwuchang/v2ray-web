package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Offered accept Offered
var Offered = []string{
	binding.MIMEJSON,
	binding.MIMEHTML,
	binding.MIMEXML,
	binding.MIMEYAML,
}

// IHelper gin 控制器
type IHelper interface {
	// 註冊 控制器
	Register(*gin.RouterGroup)
}

// Helper 爲所有 控制器 定義了 通用的輔助方法
type Helper struct {
}

// BindURI .
func (h Helper) BindURI(c *gin.Context, obj interface{}) (e error) {
	e = c.ShouldBindUri(obj)
	if e != nil {
		h.NegotiateError(c, http.StatusBadRequest, e)
		return
	}
	return
}

// NegotiateError .
func (h Helper) NegotiateError(c *gin.Context, code int, e error) {
	c.String(code, e.Error())
}

// NegotiateErrorString .
func (h Helper) NegotiateErrorString(c *gin.Context, code int, e string) {
	c.String(code, e)
}

// NegotiateData .
func (h Helper) NegotiateData(c *gin.Context, code int, data interface{}) {
	c.Negotiate(code, gin.Negotiate{
		Offered: Offered,
		Data:    data,
	})
}
