package web

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"text/template"

	"github.com/gin-gonic/gin"
	"gitlab.com/king011/v2ray-web/db/data"
)

func renderCommand(c *gin.Context, shell, text string) {
	var bufferError bytes.Buffer
	var bufferOut bytes.Buffer
	buffer := bytes.NewBufferString(text)
	cmd := exec.Command(shell)
	cmd.Stdin = buffer
	cmd.Stdout = &bufferOut
	cmd.Stderr = &bufferError
	e := cmd.Run()
	if e != nil {
		if bufferError.Len() != 0 {
			e = errors.New(bufferError.String())
		}
		c.String(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, bufferOut.String())
	return
}
func getTemplate(name string, outbound *data.Outbound, text string) (result string, e error) {
	t := template.New(name)
	t, e = t.Parse(text)
	if e != nil {
		return
	}
	ctx, e := outbound.ToContext()
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
