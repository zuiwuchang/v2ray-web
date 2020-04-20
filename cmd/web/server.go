package web

import (
	"fmt"
	"net"
	"net/http"
)

// Server 服務器
type Server struct {
	l net.Listener
}

// NewServer 創建 服務器
func NewServer(l net.Listener) (server *Server, e error) {
	server = &Server{
		l: l,
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
	fmt.Println(request.URL.Path)
}
