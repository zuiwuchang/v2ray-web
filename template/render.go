package template

import (
	"strings"
)

func Render(text string, data interface{}) (s string, e error) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "/* es5") || strings.HasPrefix(text, "// es5") {
		return renderES5(text, data)
	}
	return renderText(text, data)
}
