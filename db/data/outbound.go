package data

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"

	"gitlab.com/king011/v2ray-web/logger"

	"gitlab.com/king011/v2ray-web/utils"

	"gitlab.com/king011/v2ray-web/subscription"
)

// Outbound 可用的 出棧 配置
type Outbound struct {

	// 給人類看的 名稱
	Name string `json:"name,omitempty"`

	// 連接地址
	Add string `json:"add,omitempty"`
	// 連接端口
	Port string `json:"port,omitempty"`
	// 連接主機名
	Host string `json:"host,omitempty"`

	// 加密方案
	TLS string `json:"tls,omitempty"`

	// 使用的網路協議
	Net string `json:"net,omitempty"`

	// websocket 請求路徑
	Path string `json:"path,omitempty"`

	// 用戶身份識別碼
	UserID string `json:"userID,omitempty"`
	// 另外一個可選的用戶id
	AlterID string `json:"alterID,omitempty"`
	// Security 加密方式
	Security string `json:"security,omitempty"`
	// 用戶等級
	Level string `json:"level,omitempty"`

	// 是否是 vless 協議
	Vless bool `json:"vless,omitempty"`
}

// ToContext .
func (o *Outbound) ToContext() (context *OutboundContext, e error) {
	// vnext
	port, e := strconv.ParseUint(strings.TrimSpace(o.Port), 10, 32)
	if e != nil {
		return
	}
	aid, e := strconv.ParseInt(strings.TrimSpace(o.AlterID), 10, 64)
	level, e := strconv.ParseInt(strings.TrimSpace(o.Level), 10, 64)
	var vnext interface{}
	if o.Vless {
		vnext = gin.H{
			"address": o.Add,
			"port":    port,
			"users": []interface{}{
				gin.H{
					"id":         o.UserID,
					"encryption": o.Security,
					"level":      level,
				},
			},
		}
	} else {
		vnext = &subscription.Vnext{
			Address: o.Add,
			Port:    uint32(port),
			Users: []subscription.User{
				subscription.User{
					ID:       o.UserID,
					AlterID:  aid,
					Security: o.Security,
					Level:    level,
				},
			},
		}
	}
	// streamSettings
	streamSettings := &subscription.StreamSettings{
		Network:  o.Net,
		Security: o.TLS,
	}
	if o.TLS == "tls" {
		streamSettings.TLSSettings = &subscription.TLSSettings{
			AllowInsecure: true,
		}
	}
	if o.Net == "ws" {
		streamSettings.WebsocketSettings = &subscription.WebsocketSettings{
			Path:    o.Path,
			Headers: make(map[string]string),
		}
		if o.Host != "" {
			streamSettings.WebsocketSettings.Headers["Host"] = o.Host
		}
	}
	if o.Net == "http" || o.Net == "h2" {
		streamSettings.HTTPSettings = &subscription.HTTPSettings{
			Path: o.Path,
		}
		if o.Host != "" {
			streamSettings.HTTPSettings.Host = append(streamSettings.HTTPSettings.Host, o.Host)
		}
	}
	vnextBytes, e := json.Marshal(vnext)
	if e != nil {
		return
	}
	streamSettingsBytes, e := json.Marshal(&streamSettings)
	if e != nil {
		return
	}
	ip, e := net.LookupIP(o.Add)
	var str string
	if e == nil {
		if len(ip) > 0 {
			str = fmt.Sprint(ip[0])
		}
	} else {
		if ce := logger.Logger.Check(zap.WarnLevel, "LookupIP error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
	}
	protocol := `vmess`
	if o.Vless {
		protocol = `vless`
	}
	context = &OutboundContext{
		Outbound:       o,
		Protocol:       protocol,
		Vnext:          utils.BytesToString(vnextBytes),
		StreamSettings: utils.BytesToString(streamSettingsBytes),
		AddIP:          str,
		BasePath:       utils.BasePath(),
		Vless:          o.Vless,
	}
	return
}

// OutboundContext 模板 環境
type OutboundContext struct {
	Outbound       *Outbound
	Protocol       string
	Vnext          string
	StreamSettings string
	AddIP          string
	BasePath       string
	Vless          bool
}
