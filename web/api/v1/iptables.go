package v1

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/web"
)

// IPTables 防火牆 設定
type IPTables struct {
	web.Helper
}

// Register impl IHelper
func (h IPTables) Register(router *gin.RouterGroup) {
	r := router.Group(`iptables`)
	r.Use(h.CheckSession)

	r.GET(`view`, h.view)
	r.GET(``, h.get)
	r.GET(`default`, h.def)
	r.PUT(``, h.put)
	r.POST(`restore`, h.restore)
	r.POST(`init`, h.init)
}
func (h IPTables) renderCommand(c *gin.Context, shell, text string) {
	var bufferError bytes.Buffer
	var bufferOut bytes.Buffer
	buffer := bytes.NewBufferString(text)
	cmd := exec.Command(shell)
	cmd.Stdin = buffer
	cmd.Stdout = &bufferOut
	cmd.Stderr = &bufferError
	e := cmd.Run()
	if e != nil {
		if bufferError.Len() != 0 {
			e = errors.New(bufferError.String())
		}
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, bufferOut.String())

}

func (h IPTables) view(c *gin.Context) {
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	if strings.TrimSpace(iptables.View) == "" {
		return
	}
	h.renderCommand(c, iptables.Shell, iptables.View)
}
func (h IPTables) get(c *gin.Context) {
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, iptables)
}
func (h IPTables) def(c *gin.Context) {
	var iptables data.IPTables
	iptables.ResetDefault()
	h.NegotiateData(c, http.StatusOK, iptables)
}
func (h IPTables) put(c *gin.Context) {
	var obj data.IPTables
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.PutIPtables(&obj)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h IPTables) restore(c *gin.Context) {
	var obj data.Outbound
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	if strings.TrimSpace(iptables.Clear) == "" {
		h.NegotiateErrorString(c, http.StatusInternalServerError, `clear command nil`)
		return
	}
	h.renderCommand(c, iptables.Shell, iptables.Clear)
}
func (h IPTables) init(c *gin.Context) {
	var obj data.Outbound
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	if strings.TrimSpace(iptables.Init) == "" {
		h.NegotiateErrorString(c, http.StatusInternalServerError, `init command nil`)
		return
	}
	if strings.TrimSpace(iptables.Clear) != "" {
		var bufferError bytes.Buffer
		var bufferOut bytes.Buffer
		buffer := bytes.NewBufferString(iptables.Clear)
		cmd := exec.Command(iptables.Shell)
		cmd.Stdin = buffer
		cmd.Stdout = &bufferOut
		cmd.Stderr = &bufferError
		e = cmd.Run()
		if e != nil {
			if bufferError.Len() != 0 {
				e = errors.New(bufferError.String())
			}
			h.NegotiateError(c, http.StatusInternalServerError, e)
			return
		}
	}
	text, e := obj.ToTemplate("init", iptables.Init)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.renderCommand(c, iptables.Shell, text)
}
