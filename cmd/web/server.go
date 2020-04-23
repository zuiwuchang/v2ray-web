package web

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"text/template"

	"gitlab.com/king011/v2ray-web/cookie"

	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
	"v2ray.com/core"
	"v2ray.com/ext/tools/conf/serial"
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
		"/api/app/restore":               s.restore,
		"/api/app/login":                 s.login,
		"/api/app/logout":                s.logout,
		"/api/user/list":                 s.userList,
		"/api/user/add":                  s.userAdd,
		"/api/user/remove":               s.userRemove,
		"/api/user/password":             s.userPassword,
		"/api/v2ray/settings/get":        s.v2rayGet,
		"/api/v2ray/settings/put":        s.v2rayPut,
		"/api/v2ray/settings/test":       s.v2rayTest,
		"/api/v2ray/subscription/list":   s.subscriptionList,
		"/api/v2ray/subscription/put":    s.subscriptionPut,
		"/api/v2ray/subscription/add":    s.subscriptionAdd,
		"/api/v2ray/subscription/remove": s.subscriptionRemove,
		"/api/element/list":              s.elementList,
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
func (s *Server) getSession(helper Helper) (session *cookie.Session, e error) {
	c, e := helper.request.Cookie(cookie.CookieName)
	if e != nil {
		if e == http.ErrNoCookie {
			e = nil
		}
		return
	}
	session, e = cookie.FromCookie(c.Value)
	return
}
func (s *Server) checkSession(helper Helper) (e error) {
	c, e := helper.request.Cookie(cookie.CookieName)
	if e != nil {
		return
	}
	session, e := cookie.FromCookie(c.Value)
	if e != nil {
		return
	}
	if !session.Root {
		e = errors.New("Permission denied")
		return
	}
	return
}
func (s *Server) restore(helper Helper) (e error) {
	session, e := s.getSession(helper)
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
	c := &http.Cookie{
		Path:     "/",
		Name:     cookie.CookieName,
		Value:    val,
		HttpOnly: true,
	}
	if params.Remember {
		c.MaxAge = int(cookie.MaxAge())
	}
	http.SetCookie(helper.response, c)

	helper.RenderJSON(&session)
	return
}
func (s *Server) logout(helper Helper) (e error) {
	http.SetCookie(helper.response, &http.Cookie{
		Path:     "/",
		Name:     cookie.CookieName,
		Value:    "expired",
		MaxAge:   -1,
		HttpOnly: true,
	})
	return
}
func (s *Server) userList(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}

	var mUser manipulator.User
	result, e := mUser.List()
	if e != nil {
		return
	}
	helper.RenderJSON(result)
	return
}
func (s *Server) userAdd(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}

	var params struct {
		Name     string
		Password string
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mUser manipulator.User
	e = mUser.Add(params.Name, params.Password)
	return
}
func (s *Server) userRemove(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}

	var params struct {
		Name string
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mUser manipulator.User
	e = mUser.Remove(params.Name)
	return
}
func (s *Server) userPassword(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}

	var params struct {
		Name     string
		Password string
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mUser manipulator.User
	e = mUser.Password(params.Name, params.Password)
	return
}
func (s *Server) v2rayGet(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	text, e := mSettings.GetV2ray()
	if e != nil {
		return
	}
	helper.RenderJSON(text)
	return
}
func (s *Server) v2rayPut(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params struct {
		Text string
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.PutV2ray(params.Text)
	if e != nil {
		return
	}
	return
}
func (s *Server) v2rayTest(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params struct {
		Text string
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	outbound := &data.Outbound{
		Name:     "測試",
		Add:      "127.0.0.1",
		Port:     "1989",
		Net:      "tcp",
		Security: "auto",
		UserID:   "83b81e69-b1c7-077f-10d9-75b015b24651",
	}
	t := template.New("v2ray")
	t, e = t.Parse(params.Text)
	if e != nil {
		return
	}
	ctx, e := outbound.ToContext()
	if e != nil {
		return
	}
	var buffer bytes.Buffer
	e = t.Execute(&buffer, ctx)
	if e != nil {
		return
	}
	// v2ray
	cnf, e := serial.LoadJSONConfig(&buffer)
	if e != nil {
		return
	}
	server, e := core.New(cnf)
	if e != nil {
		return
	}
	server.Close()
	return
}
func (s *Server) subscriptionList(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var mSubscription manipulator.Subscription
	result, e := mSubscription.List()
	if e != nil {
		return
	}
	helper.RenderJSON(result)
	return
}

func (s *Server) subscriptionPut(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params data.Subscription
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}

	var mSubscription manipulator.Subscription
	e = mSubscription.Put(&params)
	if e != nil {
		return
	}
	return
}
func (s *Server) subscriptionAdd(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params data.Subscription
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mSubscription manipulator.Subscription
	e = mSubscription.Add(&params)
	if e != nil {
		return
	}
	helper.RenderJSON(params.ID)
	return
}
func (s *Server) subscriptionRemove(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params struct {
		ID uint64
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mSubscription manipulator.Subscription
	e = mSubscription.Remove(params.ID)
	if e != nil {
		return
	}
	helper.RenderJSON(params.ID)
	return
}
func (s *Server) elementList(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	element, subscription, e := mElement.List()
	if e != nil {
		return
	}

	helper.RenderJSON(&struct {
		Element      []*data.Element
		Subscription []*data.Subscription
	}{
		element,
		subscription,
	})
	return
}
