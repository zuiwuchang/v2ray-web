package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/web"
)

// Users 用戶管理
type Users struct {
	web.Helper
}

// Register impl IHelper
func (h Users) Register(router *gin.RouterGroup) {
	r := router.Group(`/users`)
	r.Use(h.CheckSession)
	r.GET(``, h.list)
	r.POST(``, h.add)
	r.PATCH(`password`, h.password)
	r.DELETE(``, h.remove)
}
func (h Users) list(c *gin.Context) {
	var mUser manipulator.User
	result, e := mUser.List()
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	h.NegotiateData(c, http.StatusOK, result)
}
func (h Users) add(c *gin.Context) {
	var obj struct {
		Name     string `form:"name" json:"name" xml:"name" yaml:"name" binding:"required"`
		Password string `form:"password" json:"password" xml:"password" yaml:"password" binding:"required"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mUser manipulator.User
	e = mUser.Add(obj.Name, obj.Password)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusCreated)
}
func (h Users) password(c *gin.Context) {
	var obj struct {
		Name     string `form:"name" json:"name" xml:"name" yaml:"name" binding:"required"`
		Password string `form:"password" json:"password" xml:"password" yaml:"password" binding:"required"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mUser manipulator.User
	e = mUser.Password(obj.Name, obj.Password)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Users) remove(c *gin.Context) {
	var obj struct {
		Name string `form:"name" json:"name" xml:"name" yaml:"name" binding:"required"`
	}
	e := h.Bind(c, &obj)
	if e != nil {
		return
	}
	var mUser manipulator.User
	e = mUser.Remove(obj.Name)
	if e != nil {
		h.NegotiateError(c, http.StatusInternalServerError, e)
		return
	}
	c.Status(http.StatusNoContent)
}
