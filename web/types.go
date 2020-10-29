package web

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// IHelper gin 控制器
type IHelper interface {
	// 註冊 控制器
	Register(*gin.RouterGroup)
}
