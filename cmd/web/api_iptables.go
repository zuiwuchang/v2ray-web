package web

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
)

type _apiIPTables struct {
}

func (a *_apiIPTables) Init(router *gin.RouterGroup) {
	r := router.Group(`/iptables`)
	GetPost(r, `/view`, a.view)
	GetPost(r, `/get`, a.get)
	GetPost(r, `/get/default`, a.getDefault)
	GetPost(r, `/put`, a.put)
	GetPost(r, `/restore`, a.restore)
	GetPost(r, `/init`, a.init)
}
func (a *_apiIPTables) view(c *gin.Context) {
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	if strings.TrimSpace(iptables.View) == "" {
		return
	}
	renderCommand(c, iptables.Shell, iptables.View)
	return
}
func (a *_apiIPTables) get(c *gin.Context) {
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, iptables)
	return
}
func (a *_apiIPTables) getDefault(c *gin.Context) {
	var iptables data.IPTables
	iptables.ResetDefault()
	c.JSON(http.StatusOK, &iptables)
	return
}
func (a *_apiIPTables) put(c *gin.Context) {
	var params data.IPTables
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSettings manipulator.Settings
	e = mSettings.PutIPtables(&params)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	return
}
func (a *_apiIPTables) restore(c *gin.Context) {
	var params data.Outbound
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	if strings.TrimSpace(iptables.Clear) == "" {
		c.String(http.StatusInternalServerError, "clear command nil")
		return
	}
	renderCommand(c, iptables.Shell, iptables.Clear)
	return
}
func (a *_apiIPTables) init(c *gin.Context) {
	var params data.Outbound
	e := c.ShouldBindWith(&params, binding.JSON)
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	var mSettings manipulator.Settings
	iptables, e := mSettings.GetIPtables()
	if e != nil {
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	if strings.TrimSpace(iptables.Init) == "" {
		c.String(http.StatusInternalServerError, `init command nil`)
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
			c.String(http.StatusInternalServerError, e.Error())
			return
		}
	}
	text, e := getTemplate("init", &params, iptables.Init)
	if e != nil {
		return
	}
	renderCommand(c, iptables.Shell, text)
	return
}
