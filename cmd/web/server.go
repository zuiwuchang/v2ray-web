package web

import (
	"bytes"
	"net"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/web/view"
	"go.uber.org/zap"
)

// Server 服務器
type Server struct {
	l      net.Listener
	router *gin.Engine
}

// NewServer 創建 服務器
func NewServer(l net.Listener) (server *Server, e error) {
	router := gin.Default()
	server = &Server{
		l:      l,
		router: router,
	}
	v := view.Helper{}
	v.Register(&router.RouterGroup)
	ms := []IModule{
		&_apiAPP{},
		&_apiIPTables{},
		&_apiProxy{},
		&_apiSettings{},
		&_apiUser{},
		&_apiV2ray{},
		&_apiWebsocket{},
	}
	r := router.Group(`api`)
	r.Use(checkRequest)
	for _, m := range ms {
		m.Init(r)
	}
	return
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
	text, e := getTemplate("init", element.Outbound, iptables.Init)
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
	return http.Serve(s.l, s.router)
}

// ServeTLS .
func (s *Server) ServeTLS(certFile, keyFile string) error {
	return http.ServeTLS(s.l,
		s.router,
		certFile, keyFile,
	)
}
