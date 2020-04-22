package data

import (
	"encoding/gob"
	"encoding/json"
	"strconv"
	"strings"

	"gitlab.com/king011/v2ray-web/utils"

	"gitlab.com/king011/v2ray-web/subscription"
)

func init() {
	gob.Register(Outbound{})
}

// Outbound 可用的 出棧 配置
type Outbound struct {
	// 數據庫記錄
	ID int64 `json:"-"`

	// 給人類看的 名稱
	Name string `json:"ps"`

	// 連接地址
	Add string `json:"add"`
	// 連接端口
	Port string `json:"port"`
	// 連接主機名
	Host string `json:"host"`

	// 加密方案
	TLS string `json:"tls"`

	// 使用的網路協議
	Net string `json:"net"`

	// websocket 請求路徑
	Path string `json:"path"`

	// 用戶身份識別碼
	UserID string `json:"id"`
	// 另外一個可選的用戶id
	AlterID string `json:"AID"`
	// Security 加密方式
	Security string `json:"type"`
	// 用戶等級
	Level string `json:"v"`
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
	vnext := &subscription.Vnext{
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
	// streamSettings
	streamSettings := &subscription.StreamSettings{
		Network:  o.Net,
		Security: o.TLS,
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
	vnextBytes, e := json.Marshal(vnext)
	if e != nil {
		return
	}
	streamSettingsBytes, e := json.Marshal(&streamSettings)
	if e != nil {
		return
	}
	context = &OutboundContext{
		Outbound:       o,
		Vnext:          utils.BytesToString(vnextBytes),
		StreamSettings: utils.BytesToString(streamSettingsBytes),
	}
	return
}

// OutboundContext 模板 環境
type OutboundContext struct {
	Outbound       *Outbound
	Vnext          string
	StreamSettings string
}
