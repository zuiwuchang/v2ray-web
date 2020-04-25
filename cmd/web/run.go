package web

import (
	"net"
	"os"

	"go.uber.org/zap"

	"gitlab.com/king011/v2ray-web/configure"
	"gitlab.com/king011/v2ray-web/logger"
)

// Run .
func Run(cnf *configure.Configure) {
	l, e := net.Listen("tcp", cnf.HTTP.Addr)
	if e == nil {
		if ce := logger.Logger.Check(zap.InfoLevel, "listen success"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.HTTP.Addr),
			)
		}
	}
	defer l.Close()
	server, e := NewServer(l, cnf.HTTP.View)
	if e != nil {
		os.Exit(1)
	}
	server.onStart()

	if cnf.HTTP.Safe() {
		if ce := logger.Logger.Check(zap.InfoLevel, "https serve"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.HTTP.Addr),
			)
		}
		e = server.ServeTLS(cnf.HTTP.CertFile, cnf.HTTP.KeyFile)
	} else {
		if ce := logger.Logger.Check(zap.InfoLevel, "http serve"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.HTTP.Addr),
			)
		}
		e = server.Serve()
	}
	if e != nil {
		if ce := logger.Logger.Check(zap.FatalLevel, "serve error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
	}
}
