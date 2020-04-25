package web

import (
	"net/http"
	"strings"
)

func (s *Server) view(response http.ResponseWriter, request *http.Request) {
	route := request.URL.Path
	filename, ok := s.m[route]
	if ok {
		s.viewFile(response, request, filename)
	} else if strings.HasPrefix(route, "/angular/zh-Hant/") {
		s.viewFile(response, request, s.root+`/zh-Hant/index.html`)
	} else if strings.HasPrefix(route, "/angular/zh-Hans/") {
		s.viewFile(response, request, s.root+`/zh-Hans/index.html`)
	} else {
		s.redirect(response, request)
	}
}
func (s *Server) redirect(response http.ResponseWriter, request *http.Request) {
	str := strings.ToLower(strings.TrimSpace(request.Header.Get("Accept-Language")))
	strs := strings.Split(str, ";")
	str = strings.TrimSpace(strs[0])
	strs = strings.Split(str, ",")
	str = strings.TrimSpace(strs[0])
	if strings.HasPrefix(str, "zh-") {
		if strings.Index(str, "cn") != -1 || strings.Index(str, "hans") != -1 {
			http.Redirect(response, request, "/angular/zh-Hans/", http.StatusFound)
		}
	}
	http.Redirect(response, request, "/angular/zh-Hant/", http.StatusFound)
}
func (s *Server) viewFile(response http.ResponseWriter, request *http.Request, filename string) {
	http.ServeFile(response, request, filename)
}
