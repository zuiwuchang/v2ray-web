package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
	"text/template"

	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/speed"
	"gitlab.com/king011/v2ray-web/utils"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
	"v2ray.com/core"
	"v2ray.com/ext/tools/conf/serial"
)

type handlerFunc func(Helper) error

// Server 服務器
type Server struct {
	l    net.Listener
	apis map[string]handlerFunc
	ws   map[string]websocket.Handler
}

// NewServer 創建 服務器
func NewServer(l net.Listener) (server *Server, e error) {
	server = &Server{
		l: l,
	}
	server.setAPI()
	server.setWebsocket()
	return
}
func (s *Server) setWebsocket() {
	s.ws = map[string]websocket.Handler{
		"/api/ws/proxy/test":   websocket.Handler(s.proxyTest),
		"/api/ws/proxy/status": websocket.Handler(s.proxyStatus),
	}
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
		"/api/proxy/list":                s.proxyList,
		"/api/proxy/update":              s.proxyUpdate,
		"/api/proxy/add":                 s.proxyAdd,
		"/api/proxy/put":                 s.proxyPut,
		"/api/proxy/remove":              s.proxyRemove,
		"/api/proxy/clear":               s.proxyClear,
		"/api/proxy/start":               s.proxyStart,
		"/api/proxy/stop":                s.proxyStop,
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
	route := request.URL.Path
	wsHandler := s.ws[route]
	if wsHandler != nil {
		wsHandler.ServeHTTP(response, request)
		return
	}
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
func (s *Server) proxyList(helper Helper) (e error) {
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
		Element      []*data.Element      `json:"element,omitempty"`
		Subscription []*data.Subscription `json:"subscription,omitempty"`
	}{
		element,
		subscription,
	})
	return
}
func (s *Server) proxyUpdate(helper Helper) (e error) {
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
	info, e := mSubscription.Get(params.ID)
	if e != nil {
		return
	}
	outbounds, e := requestSubscription(info.URL)
	if e != nil {
		return
	}
	count := len(outbounds)
	if count == 0 {
		e = errors.New("outbound empty")
		return
	}
	var mElement manipulator.Element
	result, e := mElement.Puts(info.ID, outbounds)
	if e != nil {
		return
	}
	helper.RenderJSON(result)
	return
}
func (s *Server) proxyAdd(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params struct {
		Subscription uint64
		Outbound     data.Outbound
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	result, e := mElement.Add(params.Subscription, &params.Outbound)
	if e != nil {
		return
	}
	helper.RenderJSON(result)
	return
}
func (s *Server) proxyPut(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params struct {
		ID           uint64
		Subscription uint64
		Outbound     data.Outbound
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	e = mElement.Put(params.Subscription, params.ID, &params.Outbound)
	if e != nil {
		return
	}
	return
}
func (s *Server) proxyRemove(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params struct {
		ID           uint64
		Subscription uint64
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	e = mElement.Remove(params.Subscription, params.ID)
	if e != nil {
		return
	}
	return
}
func (s *Server) proxyClear(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params struct {
		Subscription uint64
	}
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mElement manipulator.Element
	e = mElement.Clear(params.Subscription)
	if e != nil {
		return
	}
	return
}
func (s *Server) proxyStart(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params data.Element
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	e = srv.Start(&params)
	return
}
func (s *Server) proxyStop(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	srv.Stop()
	return
}
func (s *Server) proxyTest(ws *websocket.Conn) {
	defer ws.Close()
	ctx := speed.New()
	go ctx.Run()
	go func() {
		var e error
		var msg string
		for {
			if e = websocket.Message.Receive(ws, &msg); e != nil {
				ctx.Close()
				break
			}
			if msg == "close" {
				ctx.CloseSend()
				break
			}
			var element speed.Element
			e = json.Unmarshal(utils.StringToBytes(msg), &element)
			if e == nil {
				if !ctx.Send(&element) {
					break
				}
			} else {
				if ce := logger.Logger.Check(zap.WarnLevel, "unmarshal test element error"); ce != nil {
					ce.Write(
						zap.Error(e),
						zap.String("msg", msg),
					)
				}
			}
		}
	}()
	for {
		result := ctx.Get()
		if result == nil {
			break
		}
		b, e := json.Marshal(result)
		if e != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "marshal test result error"); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
		}
		e = websocket.Message.Send(ws, utils.BytesToString(b))
		if e != nil {
			ctx.Close()
			if e != io.EOF {
				if ce := logger.Logger.Check(zap.WarnLevel, "websocket closed"); ce != nil {
					ce.Write(
						zap.Error(e),
					)
				}
			}
			break
		}
	}
	websocket.Message.Send(ws, "close")
}
func (s *Server) proxyStatus(ws *websocket.Conn) {
	defer ws.Close()
	cancel := make(chan struct{})
	var closed int32
	ch := make(chan *ListenerStatus, 5)
	srv.AddListener(func(status *ListenerStatus) {
		select {
		case ch <- status:
		case <-cancel:
		}
	})
	go func() {
		var msg string
		var e error
		for {
			if e = websocket.Message.Receive(ws, &msg); e != nil {
				break
			}
		}
		if atomic.CompareAndSwapInt32(&closed, 0, 1) {
			close(cancel)
		}
	}()
	running := true
	for running {
		select {
		case <-cancel:
			running = false
		case status := <-ch:
			b, e := json.Marshal(status)
			if e == nil {
				e = websocket.Message.Send(ws, utils.BytesToString(b))
				if e != nil {
					if atomic.CompareAndSwapInt32(&closed, 0, 1) {
						close(cancel)
					}
					running = false
				}
			} else {
				if ce := logger.Logger.Check(zap.WarnLevel, "status marshal error"); ce != nil {
					ce.Write(
						zap.Error(e),
					)
				}
			}
		}
	}
}
