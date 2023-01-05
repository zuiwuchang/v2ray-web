package template

import (
	"bytes"
	"text/template"
)

func renderText(text string, data interface{}) (s string, e error) {
	t := template.New("v2ray")
	t, e = t.Parse(text)
	if e != nil {
		return
	}
	var buffer bytes.Buffer
	e = t.Execute(&buffer, data)
	if e != nil {
		return
	}
	s = buffer.String()
	return
}
