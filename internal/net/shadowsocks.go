package net

import (
	"encoding/base64"
	"net/url"
	"strings"

	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
)

func splite2(str, sep string) (first, second string, ok bool) {
	index := strings.Index(str, sep)
	if index == -1 {
		return
	}
	first = str[:index]
	second = str[index+len(sep):]
	ok = true
	return
}

type analyzeSS struct {
}

func (a *analyzeSS) do(str string) (result *Outbound) {
	str = str[len("ss://"):]
	str = strings.TrimSpace(str)
	security, userID, str, ok := a.analyzeSafe(str)
	if !ok {
		return
	}
	addr, str, ok := splite2(str, ":")
	if !ok {
		return
	}
	port, str, ok := splite2(str, "#")
	if !ok {
		return
	}
	name, e := url.QueryUnescape(str)
	if e == nil {
		str = name
	}
	result = &Outbound{
		Name:     str,
		Add:      addr,
		Port:     port,
		Security: security,
		UserID:   userID,
	}
	return
}

func (a *analyzeSS) analyzeSafe(str string) (security, userID, text string, ok bool) {
	str, text, yes := splite2(str, "@")
	if !yes {
		return
	}
	str = strings.ReplaceAll(str, "=", "")
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
	security, userID, ok = splite2(string(b), ":")
	return
}
