package web

import (
	"net/http"
	"os"
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
	if strings.Index(request.Header.Get("Accept-Encoding"), "gzip") != -1 {
		s.serveGZIP(response, request, filename)
	} else {
		http.ServeFile(response, request, filename)
	}
}

func (s *Server) toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
func (s *Server) Error(response http.ResponseWriter, e error) {
	msg, code := s.toHTTPError(e)
	http.Error(response, msg, code)
}
func (s *Server) serveGZIP(response http.ResponseWriter, request *http.Request, filename string) {
	f, e := os.Open(filename)
	if e != nil {
		s.Error(response, e)
		return
	}
	defer f.Close()
	d, e := f.Stat()
	if e != nil {
		s.Error(response, e)
		return
	}
	if d.IsDir() {
		http.Error(response, "403 Forbidden", http.StatusForbidden)
		return
	}

	size := d.Size()
	if size <= 1024 && size >= 1024*1024*5 {
		// 小於 1k 或者 大於 5m的 不壓縮
		http.ServeContent(response, request, d.Name(), d.ModTime(), f)
		return
	}

	return
}
