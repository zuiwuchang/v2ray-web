package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"text/template"

	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/speed"
	"gitlab.com/king011/v2ray-web/utils"
	"gitlab.com/king011/v2ray-web/version"
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
	m    map[string]string
	root string
}

// NewServer 創建 服務器
func NewServer(l net.Listener, root string) (server *Server, e error) {
	server = &Server{
		l:    l,
		root: root,
	}
	m := make(map[string]string)
	count := len(root)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) (e error) {
		if info.IsDir() {
			return
		}
		route := "/angular" + path[count:]
		m[route] = path
		return
	})
	server.m = m
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
		"/api/app/version":               s.version,
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
		"/api/proxy/test":                s.proxyTestOne,
		"/api/iptables/view":             s.iptablesView,
		"/api/iptables/get":              s.iptablesGet,
		"/api/iptables/get/default":      s.iptablesGetDefault,
		"/api/iptables/put":              s.iptablesPut,
		"/api/iptables/restore":          s.iptablesRestore,
		"/api/iptables/init":             s.iptablesInit,
		"/api/settings/get":              s.settingsGet,
		"/api/settings/put":              s.settingsPut,
	}
}

func (s *Server) onStart() {
	var mSettings manipulator.Settings
	result, e := mSettings.Get()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "get settings error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	} else if result == nil || !result.V2ray {
		return
	}

	element, e := mSettings.GetLast()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "get last v2ray-core error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	var iptables *data.IPTables
	if result.IPTables {
		var mSettings manipulator.Settings
		iptables, e = mSettings.GetIPtables()
		if e != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "get iptables error"); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
		}
	}
	if iptables != nil {
		s.clearIPTables(iptables)
	}
	e = srv.Start(element)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "auto start v2ray-core error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	if iptables != nil {
		s.setIPTables(iptables, element)
	}
}
func (s *Server) clearIPTables(iptables *data.IPTables) {
	if strings.TrimSpace(iptables.Init) == "" || strings.TrimSpace(iptables.Clear) == "" {
		return
	}
	buffer := bytes.NewBufferString(iptables.Clear)
	cmd := exec.Command(iptables.Shell)
	cmd.Stdin = buffer
	e := cmd.Run()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "iptables clear error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
}
func (s *Server) setIPTables(iptables *data.IPTables, element *data.Element) {
	if strings.TrimSpace(iptables.Init) == "" {
		return
	}
	text, e := s.getTemplate("init", element.Outbound, iptables.Init)
	if e != nil {
		return
	}
	buffer := bytes.NewBufferString(text)
	cmd := exec.Command(iptables.Shell)
	cmd.Stdin = buffer
	e = cmd.Run()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "iptables set error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	return
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
	if route == "/" ||
		route == "/index.html" ||
		route == "/angular/" {
		s.redirect(response, request)
		return
	}
	if strings.HasPrefix(route, "/angular/") {
		s.view(response, request)
		return
	}
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
func (s *Server) checkRequest(request *http.Request) (e error) {
	c, e := request.Cookie(cookie.CookieName)
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
func (s *Server) version(helper Helper) (e error) {
	helper.RenderJSON(map[string]string{
		"platform": fmt.Sprintf("%v %v %v", runtime.GOOS, runtime.GOARCH, runtime.Version()),
		"tag":      version.Tag,
		"commit":   version.Commit,
		"date":     version.Date,
	})
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
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	e0 := mSettings.PutLast(&params)
	if e0 != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "save last status error"); ce != nil {
			ce.Write(
				zap.Error(e0),
			)
		}
	}
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
	e := s.checkRequest(ws.Request())
	if e != nil {
		ws.Close()
		return
	}
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
func (s *Server) proxyTestOne(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params data.Outbound
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	duration, e := speed.TestOne(&params)
	if e != nil {
		return
	}
	helper.RenderJSON(duration.Milliseconds())
	return
}
func (s *Server) iptablesView(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		return
	}
	if strings.TrimSpace(iptables.View) == "" {
		return
	}
	e = s.renderCommand(helper, iptables.Shell, iptables.View)
	return
}

func (s *Server) iptablesGet(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		return
	}
	helper.RenderJSON(iptables)
	return
}
func (s *Server) iptablesGetDefault(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var iptables data.IPTables
	iptables.ResetDefault()
	helper.RenderJSON(&iptables)
	return
}
func (s *Server) iptablesPut(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params data.IPTables
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.PutIPtables(&params)
	if e != nil {
		return
	}
	return
}
func (s *Server) getTemplate(name string, outbound *data.Outbound, text string) (result string, e error) {
	t := template.New(name)
	t, e = t.Parse(text)
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
	result = buffer.String()
	return
}
func (s *Server) iptablesRestore(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params data.Outbound
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		return
	}
	if strings.TrimSpace(iptables.Clear) == "" {
		e = errors.New("clear command nil")
		return
	}
	e = s.renderCommand(helper, iptables.Shell, iptables.Clear)
	return
}
func (s *Server) iptablesInit(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params data.Outbound
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		return
	}
	if strings.TrimSpace(iptables.Init) == "" {
		e = errors.New("init command nil")
		return
	}
	if strings.TrimSpace(iptables.Clear) != "" {
		var bufferError bytes.Buffer
		var bufferOut bytes.Buffer
		buffer := bytes.NewBufferString(iptables.Clear)
		cmd := exec.Command(iptables.Shell)
		cmd.Stdin = buffer
		cmd.Stdout = &bufferOut
		cmd.Stderr = &bufferError
		e = cmd.Run()
		if e != nil {
			if bufferError.Len() != 0 {
				e = errors.New(bufferError.String())
			}
			return
		}
	}
	text, e := s.getTemplate("init", &params, iptables.Init)
	if e != nil {
		return
	}
	e = s.renderCommand(helper, iptables.Shell, text)
	return
}
func (s *Server) renderCommand(helper Helper, shell, text string) (e error) {
	var bufferError bytes.Buffer
	var bufferOut bytes.Buffer
	buffer := bytes.NewBufferString(text)
	cmd := exec.Command(shell)
	cmd.Stdin = buffer
	cmd.Stdout = &bufferOut
	cmd.Stderr = &bufferError
	e = cmd.Run()
	if e != nil {
		if bufferError.Len() != 0 {
			e = errors.New(bufferError.String())
		}
		return
	}
	helper.RenderJSON(bufferOut.String())
	return
}
func (s *Server) settingsGet(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	result, e := mSettings.Get()
	if e != nil {
		return
	}
	helper.RenderJSON(result)
	return
}
func (s *Server) settingsPut(helper Helper) (e error) {
	e = s.checkSession(helper)
	if e != nil {
		return
	}
	var params data.Settings
	e = helper.BodyJSON(&params)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.Put(&params)
	return
}
