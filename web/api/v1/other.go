package v1

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xtls/xray-core/core"
	"gitlab.com/king011/v2ray-web/internal/logs"
	"gitlab.com/king011/v2ray-web/utils"
	"gitlab.com/king011/v2ray-web/version"
	"gitlab.com/king011/v2ray-web/web"
)

var startAt = time.Now()

// Other 一些其它的 api
type Other struct {
	web.Helper
}

// Register impl IHelper
func (h Other) Register(router *gin.RouterGroup) {
	router.GET(`version`, h.version)

	router.GET(`logs/websocket`, h.CheckSession, h.logs)
}
func (h Other) version(c *gin.Context) {
	gv := gin.Version
	if strings.HasPrefix(gv, "v") {
		gv = gv[1:]
	}

	h.NegotiateFile(c, `version`, startAt, gin.H{
		`platform`: fmt.Sprintf(`%s %s %s gin%s`, runtime.GOOS, runtime.GOARCH, runtime.Version(), gv),
		"tag":      version.Version,
		"commit":   version.Commit,
		"date":     version.Date,
		"v2ray":    core.Version(),
	})
}
func (h Other) logs(c *gin.Context) {
	ws, e := h.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	defer ws.Close()
	cancel := make(chan struct{})
	var closed int32
	ch := make(chan string, 5)
	id := logs.AddListener(func(str string) {
		select {
		case ch <- str:
		default:
		}
	})
	go func() {
		var e error
		for {
			_, _, e = ws.ReadMessage()
			if e != nil {
				break
			}
		}
		if atomic.CompareAndSwapInt32(&closed, 0, 1) {
			close(cancel)
		}
	}()
	running := true
	for running {
		select {
		case <-cancel:
			running = false
		case str := <-ch:
			e := ws.WriteMessage(websocket.TextMessage, utils.StringToBytes(str))
			if e != nil {
				if atomic.CompareAndSwapInt32(&closed, 0, 1) {
					close(cancel)
				}
				running = false
			}
		}
	}
	logs.RemoveListener(id)
	return
}
