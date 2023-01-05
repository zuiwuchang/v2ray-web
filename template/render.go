package template

import (
	_ "embed"
	"strings"
)

//go:embed reader.js
var reader string

//go:embed reader_route.js
var route string

//go:embed reader_proxy.js
var proxy string

var Default string
var Proxy string

func init() {
	Default = reader + "\n" + route
	Proxy = reader + "\n" + proxy
}

func Render(text string, data interface{}) (s string, e error) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "/* es5") || strings.HasPrefix(text, "// es5") {
		return renderES5(text, data)
	}
	return renderText(text, data)
}
