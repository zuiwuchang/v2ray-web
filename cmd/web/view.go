package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// GzipDefaultExcludedExtentions 默認排除 擴展名
var GzipDefaultExcludedExtentions = []string{
	`.png`, `.gif`, `.jpeg`, `.jpg`,
	`.mp4`,
	`.mp3`,
	`.gz`, `.gzip`, `.zip`, `.bz`, `.bz2`, `.bzip2`, `.xz`,
	`.rar`, `.7z`,
	`.woff`, `.woff2`,
}

type _View struct {
	m      map[string]string
	root   string
	filter map[string]bool
}

func (v *_View) Init(router *gin.RouterGroup) {
	filter := make(map[string]bool, len(GzipDefaultExcludedExtentions)+1)
	for _, k := range GzipDefaultExcludedExtentions {
		filter[k] = true
	}
	filter[`.ico`] = true
	v.filter = filter
	gz := gzip.Gzip(gzip.DefaultCompression)

	root := v.root
	GetPost(router, `/`, v.redirect)
	GetPost(router, `/index.html`, v.redirect)
	GetPost(router, `/angular`, v.redirect)
	GetPost(router, `/angular/*path`, v.gzipFilter, gz, v.view)

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
	str := strings.ToLower(strings.TrimSpace(request.Header.Get(`Accept-Language`)))
	strs := strings.Split(str, `;`)
	str = strings.TrimSpace(strs[0])
	strs = strings.Split(str, `,`)
	str = strings.TrimSpace(strs[0])
	if strings.HasPrefix(str, `zh-`) {
		if strings.Index(str, `cn`) != -1 || strings.Index(str, `hans`) != -1 {
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
func (v *_View) gzipFilter(c *gin.Context) {
	str := c.Request.Header.Get(`Accept-Encoding`)
	if !strings.Contains(str, `gzip`) {
		return
	}

	path := c.Param(`path`)
	filename, ok := v.m[path]
	if !ok {
		return
	}
	ext := strings.ToLower(filepath.Ext(filename))
	if v.filter[ext] {
		c.Request.Header.Set(`Accept-Encoding`, strings.ReplaceAll(str, `gzip`, ``))
		return
	}

	stat, _ := os.Stat(filename)
	if stat != nil {
		size := stat.Size()
		// 小於 1k 大於 5m 不使用 gzip
		if size < 1024 || size > 1024*1024*5 {
			c.Request.Header.Set(`Accept-Encoding`, strings.ReplaceAll(str, `gzip`, ``))
			return
		}
	}
}
