package web

import (
	"io"
	"net/http"

	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/speed"
	"gitlab.com/king011/v2ray-web/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type _apiWebsocket struct {
	upgrader websocket.Upgrader
}

func (a *_apiWebsocket) Init(router *gin.RouterGroup) {
	a.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	r := router.Group(`/ws`)
	GetPost(r, `/proxy/test`, a.proxyTest)
}

func (a *_apiWebsocket) proxyTest(c *gin.Context) {
	ws, e := a.upgrader.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	defer ws.Close()

	ctx := speed.New(getURL())
	go ctx.Run()
	go func() {
		var e error
		var b []byte
		var t int
		for {
			t, b, e = ws.ReadMessage()
			if e != nil {
				ctx.Close()
				break
			} else if t != websocket.TextMessage {
				continue
			}
			msg := utils.BytesToString(b)
			if msg == "close" {
				ctx.CloseSend()
				break
			}
			var element speed.Element
			e = json.Unmarshal(b, &element)
			if e == nil {
				if !ctx.Send(&element) {
					break
				}
			} else {
				if ce := logger.Logger.Check(zap.WarnLevel, "unmarshal test element error"); ce != nil {
					ce.Write(
						zap.Error(e),
						zap.String("msg", msg),
					)
				}
			}
		}
	}()
	for {
		result := ctx.Get()
		if result == nil {
			break
		}
		b, e := json.Marshal(result)
		if e != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "marshal test result error"); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
		}
		e = ws.WriteMessage(websocket.TextMessage, b)
		if e != nil {
			ctx.Close()
			if e != io.EOF {
				if ce := logger.Logger.Check(zap.WarnLevel, "websocket closed"); ce != nil {
					ce.Write(
						zap.Error(e),
					)
				}
			}
			break
		}
	}
	e = ws.WriteMessage(websocket.TextMessage, []byte(`close`))
}
