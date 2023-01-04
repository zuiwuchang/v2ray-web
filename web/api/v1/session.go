package v1

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/web"
	"go.uber.org/zap"
)

var _sessionLock _SessionLock
var errSessionLock = errors.New(`lock session error`)
var errIPLock = errors.New(`lock ip error`)

func init() {
	_sessionLock.keys = make(map[string]int8)
	_sessionLock.keysIP = make(map[string]int8)
}

// 防止暴力猜密碼
type _SessionLock struct {
	keys   map[string]int8
	keysIP map[string]int8
	mutex  sync.Mutex
}

func (s *_SessionLock) Lock(key string) (e error) {
	s.mutex.Lock()
	num := s.keys[key]
	if num < 2 {
		num++
		s.keys[key] = num
	} else {
		e = errSessionLock
	}
	s.mutex.Unlock()
	return
}
func (s *_SessionLock) Unlock(key string) {
	s.mutex.Lock()
	if num, ok := s.keys[key]; ok {
		if num < 2 {
			if len(s.keys) > 30 {
				delete(s.keys, key)
			} else {
				s.keys[key] = 0
			}
		} else {
			s.keys[key] = num - 1
		}
	}
	s.mutex.Unlock()
}
func (s *_SessionLock) LockIP(key string) (e error) {
	s.mutex.Lock()
	num := s.keysIP[key]
	if num < 90 {
		num++
		s.keysIP[key] = num
	} else {
		e = errIPLock
	}
	s.mutex.Unlock()
	return
}
func (s *_SessionLock) UnlockIP(key string) {
	s.mutex.Lock()
	if num, ok := s.keysIP[key]; ok {
		if num < 2 {
			if len(s.keysIP) > 90 {
				delete(s.keysIP, key)
			} else {
				s.keysIP[key] = 0
			}
		} else {
			s.keysIP[key] = num - 1
		}
	}
	s.mutex.Unlock()
}

// Session 會話
type Session struct {
	web.Helper
}

// Register impl IHelper
func (h Session) Register(router *gin.RouterGroup) {
	r := router.Group(`session`)
	r.POST(``, h.login)
	r.GET(`:at/:maxage/:token`, h.restore)
}

func (h Session) login(c *gin.Context) {
	// 解析參數
	var obj struct {
		Name     string `form:"name" json:"name" xml:"name" yaml:"name" binding:"required"`
		Password string `form:"password" json:"password" xml:"password" yaml:"password" binding:"required"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	ip := c.ClientIP()
	e = _sessionLock.LockIP(ip)
	if e != nil {
		h.NegotiateError(c, http.StatusForbidden, e)
		return
	}
	defer _sessionLock.UnlockIP(ip)
	key := ip + obj.Name
	e = _sessionLock.Lock(key)
	if e != nil {
		h.NegotiateError(c, http.StatusForbidden, e)
		return
	}
	defer _sessionLock.Unlock(key)

	// 查詢用戶
	var mUser manipulator.User
	session, e := mUser.Login(obj.Name, obj.Password)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	} else if session == nil {
		h.NegotiateErrorString(c, http.StatusNotFound, `name or password not match`)
		return
	}
	// 生成 cookie
	token, e := session.Cookie()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}

	// 響應 數據
	h.NegotiateData(c, http.StatusCreated, gin.H{
		`session`: session,
		`token`:   token,
		`maxage`:  cookie.MaxAge(),
	})

	if ce := logger.Logger.Check(zap.WarnLevel, c.FullPath()); ce != nil {
		ce.Write(
			zap.String(`method`, c.Request.Method),
			zap.String(`session`, session.String()),
			zap.String(`client ip`, c.ClientIP()),
		)
	}
}
func (h Session) restore(c *gin.Context) {
	var obj struct {
		Token  string `uri:"token" binding:"required"`
		At     int64  `uri:"at" binding:"required"`
		Maxage int64  `uri:"maxage" binding:"required"`
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	obj.Token, e = url.QueryUnescape(obj.Token)
	if e != nil {
		h.NegotiateError(c, http.StatusBadRequest, e)
		return
	}
	session, e := cookie.FromCookie(obj.Token)
	if e != nil {
		h.NegotiateError(c, http.StatusUnauthorized, e)
		return
	}
	var last time.Time
	now := time.Now()
	if obj.Maxage > 0 && obj.At > 0 {
		last = time.Unix(obj.At, 0).Local()
		at := time.Unix(obj.At, 0).Local()
		expired := at.Add(time.Duration(obj.Maxage) * time.Second)
		maxage := expired.Sub(now) / time.Second
		if maxage < 0 {
			maxage = 0
		}
		c.Header("Cache-Control", "max-age="+strconv.Itoa((int(maxage))))
	}
	if last.IsZero() {
		last = now
	}
	h.NegotiateFile(c, `token`, last, session)
	if ce := logger.Logger.Check(zap.InfoLevel, c.FullPath()); ce != nil {
		ce.Write(
			zap.String(`method`, c.Request.Method),
			zap.String(`session`, session.String()),
			zap.String(`client ip`, c.ClientIP()),
		)
	}
}
