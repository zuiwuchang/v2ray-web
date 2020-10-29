package web

import (
	"encoding/xml"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/yaml.v2"
)

func (h Helper) toHTTPError(c *gin.Context, e error) {
	if os.IsNotExist(e) {
		h.NegotiateError(c, http.StatusNotFound, e)
		return
	}
	if os.IsPermission(e) {
		h.NegotiateError(c, http.StatusForbidden, e)
		return
	}
	h.NegotiateError(c, http.StatusInternalServerError, e)
}

// NegotiateFilesystem .
func (h Helper) NegotiateFilesystem(c *gin.Context, fs http.FileSystem, path string) {
	if path == `/` || path == `` {
		path = `/index.html`
	}
	f, e := fs.Open(path)
	if e != nil {
		if os.IsNotExist(e) {
			path = `/index.html`
			f, e = fs.Open(path)
		}
	}
	if e != nil {
		h.toHTTPError(c, e)
		return
	}
	stat, e := f.Stat()
	if e != nil {
		f.Close()
		h.toHTTPError(c, e)
		return
	}
	if stat.IsDir() {
		f.Close()
		h.NegotiateErrorString(c, http.StatusForbidden, `not a file`)
		return
	}

	_, name := filepath.Split(path)
	http.ServeContent(c.Writer, c.Request, name, stat.ModTime(), f)
	f.Close()
}

// NegotiateFile .
func (h Helper) NegotiateFile(c *gin.Context, name string, modtime time.Time, obj interface{}) {
	content := &ReadSeeker{
		Data: obj,
	}
	switch c.NegotiateFormat(Offered...) {
	case binding.MIMEXML:
		c.Header("Content-Type", "application/xml; charset=utf-8")
		content.Marshal = xml.Marshal
	case binding.MIMEYAML:
		c.Header("Content-Type", "application/x-yaml; charset=utf-8")
		content.Marshal = yaml.Marshal
	default:
		// 默認以 json
		c.Header("Content-Type", "application/json; charset=utf-8")
		content.Marshal = json.Marshal
	}
	http.ServeContent(c.Writer, c.Request, name, modtime, content)
}
