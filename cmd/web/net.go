package web

import (
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

func requestSubscription(url string) (result []*data.Outbound, e error) {
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

	dst, e := base64.RawStdEncoding.DecodeString(utils.BytesToString(src))
	if e != nil {
		return
	}
	str := utils.BytesToString(dst)
	strs := strings.Split(str, "\n")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		node := analyzeString(str)
		if node != nil {
			result = append(result, node)
		}
	}
	return
}
func analyzeString(str string) (result *data.Outbound) {
	str = strings.TrimSpace(str)
	if !strings.HasPrefix(str, "vmess://") {
		if ce := logger.Logger.Check(zap.WarnLevel, "not support outbound"); ce != nil {
			ce.Write(
				zap.String("value", str),
			)
		}
		return
	}
	str = str[len("vmess://"):]
	b, e := base64.StdEncoding.DecodeString(str)
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
	var node data.Outbound
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

var replaceNumber = regexp.MustCompile(`"\s*:\s*([\d]+)\s*,`)
