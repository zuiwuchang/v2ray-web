package v1

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/version"
	"gitlab.com/king011/v2ray-web/web"
	"v2ray.com/core"
)

var startAt = time.Now()

// Other 一些其它的 api
type Other struct {
	web.Helper
}

// Register impl IHelper
func (h Other) Register(router *gin.RouterGroup) {
	router.GET(`/version`, h.version)
}
func (h Other) version(c *gin.Context) {
	gv := gin.Version
	if strings.HasPrefix(gv, "v") {
		gv = gv[1:]
	}

	h.NegotiateFile(c, `version`, startAt, gin.H{
		`platform`: fmt.Sprintf(`%s %s %s gin%s`, runtime.GOOS, runtime.GOARCH, runtime.Version(), gv),
		"tag":      version.Tag,
		"commit":   version.Commit,
		"date":     version.Date,
		"v2ray":    core.Version(),
	})
}
