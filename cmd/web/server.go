package web

import (
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"gitlab.com/king011/v2ray-web/cookie"

	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
)

type handlerFunc func(Helper) error

// Server 服務器
type Server struct {
	l    net.Listener
	apis map[string]handlerFunc
}

// NewServer 創建 服務器
func NewServer(l net.Listener) (server *Server, e error) {
	server = &Server{
		l: l,
	}
	server.setAPI()
	return
}
func (s *Server) setAPI() {
	s.apis = map[string]handlerFunc{
		"/api/app/restore": s.restore,
		"/api/app/login":   s.login,
		"/api/app/logout":  s.logout,
	}
}

// Serve .
func (s *Server) Serve() error {
	return http.Serve(s.l, s)
}

// ServeTLS .
func (s *Server) ServeTLS(certFile, keyFile string) error {
	return http.ServeTLS(s.l,
		s,
		certFile, keyFile,
	)
}
func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var helper Helper
	if request.Body != nil {
		body, e := ioutil.ReadAll(io.LimitReader(request.Body, 1024*32))
		request.Body.Close()
		if e != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		helper = Helper{
			request:  request,
			body:     body,
			response: response,
		}
	}
	route := request.URL.Path
	handler := s.apis[route]
	if handler != nil {
		e := handler(helper)
		if e != nil {
			helper.RenderError(e)
		}
		return
	}

	if ce := logger.Logger.Check(zap.WarnLevel, "route not found"); ce != nil {
		ce.Write(
			zap.String("route", request.URL.Path),
		)
	}
	helper.RenderText(http.StatusNotFound, "route not found")
}

func (s *Server) restore(helper Helper) (e error) {
	c, e := helper.request.Cookie(cookie.CookieName)
	if e != nil {
		if e == http.ErrNoCookie {
			e = nil
		}
		return
	}
	session, e := cookie.FromCookie(c.Value)
	if e != nil {
		return
	}
	helper.RenderJSON(session)
	return
}
func (s *Server) login(helper Helper) (e error) {
	var params struct {
		Name     string
		Password string
		Remember bool
	}
	e = helper.BodyJSON(&params)
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
	http.SetCookie(helper.response, &http.Cookie{
		Name:     cookie.CookieName,
		Value:    val,
		MaxAge:   int(cookie.MaxAge()),
		HttpOnly: true,
	})
	helper.RenderJSON(&session)
	return
}
func (s *Server) logout(helper Helper) (e error) {
	http.SetCookie(helper.response, &http.Cookie{
		Name:     cookie.CookieName,
		MaxAge:   -1,
		HttpOnly: true,
	})
	return
}
