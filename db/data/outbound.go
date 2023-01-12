package data

import (
	"bytes"
	"fmt"
	"net"
	text_template "text/template"

	"gitlab.com/king011/v2ray-web/template"
	"go.uber.org/zap"

	"gitlab.com/king011/v2ray-web/logger"

	"gitlab.com/king011/v2ray-web/utils"
)

// Outbound 可用的 出棧 配置
type Outbound struct {

	// 給人類看的 名稱
	Name string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"`

	// 連接地址
	Add string `json:"add,omitempty" xml:"add,omitempty" yaml:"add,omitempty"`
	// 連接端口
	Port string `json:"port,omitempty" xml:"port,omitempty" yaml:"port,omitempty"`
	// 連接主機名
	Host string `json:"host,omitempty" xml:"host,omitempty" yaml:"host,omitempty"`

	// 加密方案
	TLS string `json:"tls,omitempty" xml:"tls,omitempty" yaml:"tls,omitempty"`

	// 使用的網路協議
	Net string `json:"net,omitempty" xml:"net,omitempty" yaml:"net,omitempty"`

	// websocket 請求路徑
	Path string `json:"path,omitempty" xml:"path,omitempty" yaml:"path,omitempty"`

	// 用戶身份識別碼
	UserID string `json:"userID,omitempty" xml:"userID,omitempty" yaml:"userID,omitempty"`
	// 另外一個可選的用戶id
	AlterID string `json:"alterID,omitempty" xml:"alterID,omitempty" yaml:"alterID,omitempty"`
	// Security 加密方式
	Security string `json:"security,omitempty" xml:"security,omitempty" yaml:"security,omitempty"`
	// 用戶等級
	Level string `json:"level,omitempty" xml:"level,omitempty" yaml:"level,omitempty"`

	// 協議 名稱
	Protocol string `json:"protocol,omitempty" xml:"protocol,omitempty" yaml:"protocol,omitempty"`
	// xtls 流控
	Flow string `json:"flow,omitempty" xml:"flow,omitempty" yaml:"flow,omitempty"`
}

// ToTemplate .
func (o *Outbound) ToTemplate(name, text string) (result string, e error) {
	t := text_template.New(name)
	t, e = t.Parse(text)
	if e != nil {
		return
	}
	ctx, e := o.ToContext()
	if e != nil {
		return
	}
	var buffer bytes.Buffer
	e = t.Execute(&buffer, ctx)
	if e != nil {
		return
	}
	result = buffer.String()
	return
}

// ToContext .
func (o *Outbound) ToContext() (context *OutboundContext, e error) {
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
	context = &OutboundContext{
		Outbound: o,
		AddIP:    str,
		BasePath: utils.BasePath(),
	}
	return
}
func (o *Outbound) Render(text string) (string, error) {
	ctx, e := o.ToContext()
	if e != nil {
		return ``, e
	}
	return template.Render(text, ctx)
}
func (o *Outbound) RenderTarget(port int) (string, error) {
	ctx, e := o.ToContext()
	if e != nil {
		return ``, e
	}
	return template.Render(template.Proxy, map[string]interface{}{
		`ctx`:  ctx,
		`port`: port,
	})
}
func (o *Outbound) RenderStrategy(text string, strategy *StrategyValue) (string, error) {
	ctx, e := o.ToContext()
	ctx.Strategy = strategy
	if e != nil {
		return ``, e
	}
	return template.Render(text, ctx)
}

// OutboundContext 模板 環境
type OutboundContext struct {
	Outbound *Outbound
	AddIP    string
	BasePath string
	Strategy *StrategyValue
}
