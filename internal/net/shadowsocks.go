package net

import (
	"encoding/base64"
	"net/url"
	"strings"

	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
)

type analyzeSS struct {
}

func (a *analyzeSS) do(str string) (result *Outbound) {
	u, e := url.Parse(str)
	if e != nil {
		return
	}
	var userID, security string
	if u.User != nil {
		user := strings.TrimRight(u.User.Username(), "=")
		b, e := base64.RawStdEncoding.DecodeString(user)
		if e != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "decode base64 outbound error"); ce != nil {
				ce.Write(
					zap.Error(e),
					zap.String("value", str),
				)
			}
			return
		}
		strs := strings.SplitN(string(b), ":", 2)
		security = strs[0]
		if len(strs) > 1 {
			userID = strs[1]
		}
	}
	result = &Outbound{
		Add:      u.Hostname(),
		Port:     u.Port(),
		Name:     u.Fragment,
		UserID:   userID,
		Security: security,
	}
	return
}
