package web

import (
	"crypto/tls"
	"net"
	"os"

	"go.uber.org/zap"

	"gitlab.com/king011/v2ray-web/configure"
	"gitlab.com/king011/v2ray-web/logger"
)

// Run .
func Run(cnf *configure.Configure) {
	l, e := listen(&cnf.HTTP)
	if e != nil {
		if ce := logger.Logger.Check(zap.FatalLevel, "listen error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		os.Exit(1)
	}
	defer l.Close()
}
func listen(cnf *configure.HTTP) (l net.Listener, e error) {
	if cnf.Safe() {
		var cert tls.Certificate
		cert, e = tls.LoadX509KeyPair(cnf.CertFile, cnf.KeyFile)
		if e != nil {
			return
		}
		l, e = tls.Listen("tcp", cnf.Addr, &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		if e == nil {
			if ce := logger.Logger.Check(zap.InfoLevel, "https listen"); ce != nil {
				ce.Write(
					zap.String("addr", cnf.Addr),
				)
			}
		}
	} else {
		l, e = net.Listen("tcp", cnf.Addr)
		if e == nil {
			if ce := logger.Logger.Check(zap.InfoLevel, "http listen"); ce != nil {
				ce.Write(
					zap.String("addr", cnf.Addr),
				)
			}
		}
	}
	return
}
