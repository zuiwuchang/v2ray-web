package net

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"gitlab.com/king011/v2ray-web/utils"

	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
)

// Outbound 可用的 出棧 配置
type Outbound struct {

	// 給人類看的 名稱
	Name string `json:"ps,omitempty"`

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
	UserID string `json:"id,omitempty"`
	// 另外一個可選的用戶id
	AlterID string `json:"aid,omitempty"`
	// Security 加密方式
	Security string `json:"type,omitempty"`
	// 用戶等級
	Level string `json:"v,omitempty"`
}

// RequestSubscription .
func RequestSubscription(url string) (result []*data.Outbound, e error) {
	response, e := http.Get(url)
	if e != nil {
		return
	}
	var src []byte
	var er error
	if response.Body != nil {
		src, er = ioutil.ReadAll(response.Body)
	}
	if response.StatusCode != 200 {
		e = fmt.Errorf("%v %v", response.StatusCode, response.Status)
		return
	}
	if er != nil {
		e = er
		return
	}
	src = bytes.TrimSpace(src)
	str := strings.ReplaceAll(string(src), "=", "")
	dst, e := base64.RawStdEncoding.DecodeString(str)
	if e != nil {
		return
	}
	str = utils.BytesToString(dst)
	strs := strings.Split(str, "\n")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		protocol, node := AnalyzeString(str)
		if node != nil {
			result = append(result, &data.Outbound{
				Name:     node.Name,
				Add:      node.Add,
				Port:     node.Port,
				Host:     node.Host,
				TLS:      node.TLS,
				Net:      node.Net,
				Path:     node.Path,
				UserID:   node.UserID,
				AlterID:  node.AlterID,
				Security: node.Security,
				Level:    node.Level,
				Protocol: protocol,
			})
		}
	}
	return
}

// AnalyzeString 分享 鏈接
func AnalyzeString(str string) (protocol string, result *Outbound) {
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "vless://") {
		protocol = "vless"
		result = analyzeVMess(str)
	} else if strings.HasPrefix(str, "vmess://") {
		protocol = "vmess"
		result = analyzeVMess(str)
	} else if strings.HasPrefix(str, "ss://") {
		protocol = "shadowsocks"
		var analyze analyzeSS
		result = analyze.do(str)
	} else {
		if ce := logger.Logger.Check(zap.WarnLevel, "not support outbound"); ce != nil {
			ce.Write(
				zap.String("value", str),
			)
		}
	}
	return
}

var replaceNumber = regexp.MustCompile(`"\s*:\s*([\d]+)\s*,`)

func analyzeVMess(str string) (result *Outbound) {
	str = str[len("vmess://"):]
	str = strings.ReplaceAll(strings.TrimSpace(str), "=", "")
	b, e := base64.RawStdEncoding.DecodeString(str)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "decode base64 outbound error"); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String("value", str),
			)
		}
		return
	}
	b = replaceNumber.ReplaceAll(b, []byte(`":"$1",`))
	var node Outbound
	e = json.Unmarshal(b, &node)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "unmarshal outbound error"); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String("value", str),
			)
		}
		return
	}
	result = &node
	return
}
