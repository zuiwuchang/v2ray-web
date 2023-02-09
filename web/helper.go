package web

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// Offered accept Offered
var Offered = []string{
	binding.MIMEJSON,
	binding.MIMEHTML,
	binding.MIMEXML,
	binding.MIMEYAML,
}
var _compression = gzip.Gzip(gzip.DefaultCompression)

// Helper 輔助類型
type Helper struct {
}

// NegotiateData .
func (h Helper) NegotiateData(c *gin.Context, code int, data interface{}) {
	switch c.NegotiateFormat(Offered...) {
	case binding.MIMEXML:
		c.XML(code, data)
	case binding.MIMEYAML:
		c.YAML(code, data)
	default:
		// 默認以 json
		c.JSON(code, data)
	}
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

// Bind .
func (h Helper) Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return h.BindWith(c, obj, b)
}

// BindWith .
func (h Helper) BindWith(c *gin.Context, obj interface{}, b binding.Binding) (e error) {
	e = c.ShouldBindWith(obj, b)
	if e != nil {
		h.NegotiateError(c, http.StatusBadRequest, e)
		return
	}
	return
}

// ShouldBindSession 返回session 不進行響應
func (h Helper) ShouldBindSession(c *gin.Context) (session *cookie.Session, e error) {
	v, exists := c.Get(`session`)
	if exists {
		if v == nil {
			return
		} else if tmp, ok := v.(error); ok {
			e = tmp
			return
		} else if tmp, ok := v.(*cookie.Session); ok {
			session = tmp
			return
		}
		if ce := logger.Logger.Check(zap.ErrorLevel, `unknow session type`); ce != nil {
			ce.Write(
				zap.String(`method`, c.Request.Method),
				zap.String(`session`, fmt.Sprint(session)),
				zap.String(`session type`, fmt.Sprint(reflect.TypeOf(session))),
			)
		}
		return
	}
	token, e := h.getToken(c)
	if e != nil {
		if ce := logger.Logger.Check(zap.ErrorLevel, `get token`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`method`, c.Request.Method),
			)
		}
		return
	}
	if token == "" {
		c.Set(`session`, nil)
		return
	}
	session, e = cookie.FromCookie(token)
	if e == nil {
		c.Set(`session`, session)
	} else {
		c.Set(`session`, e)
	}
	return
}
func (h Helper) getToken(c *gin.Context) (value string, e error) {
	value = c.GetHeader("token")
	if value != "" {
		return
	}
	if !c.IsWebsocket() {
		return
	}
	var obj struct {
		Token string `form:"token"`
	}
	e = c.ShouldBindQuery(&obj)
	if e != nil {
		return
	}
	value = obj.Token
	return
}

// BindSession 返回 session 並響應錯誤
func (h Helper) BindSession(c *gin.Context) (result *cookie.Session) {
	session, e := h.ShouldBindSession(c)
	if e != nil {
		h.NegotiateError(c, http.StatusUnauthorized, e)
		return
	} else if session == nil {
		h.NegotiateErrorString(c, http.StatusUnauthorized, `session miss`)
		return
	}
	result = session
	return
}

// CheckSession 檢查是否具有 session
func (h Helper) CheckSession(c *gin.Context) {
	session := h.BindSession(c)
	if session == nil {
		c.Abort()
		return
	}
}

// BindQuery .
func (h Helper) BindQuery(c *gin.Context, obj interface{}) error {
	return h.BindWith(c, obj, binding.Query)
}

// Compression .
func (h Helper) Compression() gin.HandlerFunc {
	return _compression
}

// Compression .
func Compression() gin.HandlerFunc {
	return _compression
}

// Upgrade .
func (h Helper) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	ws, e := upgrader.Upgrade(w, r, responseHeader)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `websocket error`); ce != nil {
			ce.Write(
				zap.String(`path`, r.URL.Path),
				zap.Error(e),
			)
		}
	}
	return ws, e
}

// JSON .
func (h Helper) JSON() jsoniter.API {
	return json
}
