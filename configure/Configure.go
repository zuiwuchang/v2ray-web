package configure

import (
	"encoding/json"
	"io/ioutil"

	"github.com/google/go-jsonnet"
	logger "gitlab.com/king011/king-go/log/logger.zap"
)

// Configure global configure
type Configure struct {
	Logger   logger.Options
	HTTP     HTTP
	Database Database
}

// Format format global configure
func (c *Configure) Format(basePath string) (e error) {
	e = c.HTTP.Format(basePath)
	if e != nil {
		return
	}
	e = c.Database.Format(basePath)
	if e != nil {
		return
	}
	return
}
func (c *Configure) String() string {
	if c == nil {
		return "nil"
	}
	b, e := json.MarshalIndent(c, "", "	")
	if e != nil {
		return e.Error()
	}
	return string(b)
}

var _Configure Configure

// Single single Configure
func Single() *Configure {
	return &_Configure
}

// Load load configure file
func (c *Configure) Load(filename string) (e error) {
	var b []byte
	b, e = ioutil.ReadFile(filename)
	if e != nil {
		return
	}
	vm := jsonnet.MakeVM()
	var jsonStr string
	jsonStr, e = vm.EvaluateSnippet("", string(b))
	if e != nil {
		return
	}
	b = []byte(jsonStr)
	e = json.Unmarshal(b, c)
	if e != nil {
		return
	}
	return
}
