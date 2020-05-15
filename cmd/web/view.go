package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type _View struct {
	m    map[string]string
	root string
}

func (v *_View) RegisterTo(router *gin.RouterGroup, root string) {
	v.root = root
	router.GET(`/`, v.redirect)
	router.GET(`/index.html`, v.redirect)
	router.GET(`/angular`, v.redirect)
	router.GET(`/angular/*path`, v.view)

	m := make(map[string]string)
	count := len(root)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) (e error) {
		if info.IsDir() {
			return
		}

		route := strings.ReplaceAll(path[count:], `\`, `/`)
		m[route] = path
		return
	})
	v.m = m
}
func (v *_View) redirect(c *gin.Context) {
	request := c.Request
	str := strings.ToLower(strings.TrimSpace(request.Header.Get("Accept-Language")))
	strs := strings.Split(str, ";")
	str = strings.TrimSpace(strs[0])
	strs = strings.Split(str, ",")
	str = strings.TrimSpace(strs[0])
	if strings.HasPrefix(str, "zh-") {
		if strings.Index(str, "cn") != -1 || strings.Index(str, "hans") != -1 {
			c.Redirect(http.StatusFound, `/angular/zh-Hans/`)
			return
		}
	}
	c.Redirect(http.StatusFound, `/angular/zh-Hant/`)
}
func (v *_View) view(c *gin.Context) {
	path := c.Param(`path`)
	filename, ok := v.m[path]
	if ok {
		c.File(filename)
	} else if strings.HasPrefix(path, "/zh-Hant/") {
		c.File(v.root + `/zh-Hant/index.html`)
	} else if strings.HasPrefix(path, "/zh-Hans/") {
		c.File(v.root + `/zh-Hans/index.html`)
	} else {
		v.redirect(c)
	}
}
