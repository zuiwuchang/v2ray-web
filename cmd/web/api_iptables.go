package web

import (
	"github.com/gin-gonic/gin"
)

type _apiIPTables struct {
}

func (a *_apiIPTables) Init(router *gin.RouterGroup) {
	r := router.Group(`/iptables`)

	GetPost(r, `/restore`, a.restore)
	GetPost(r, `/init`, a.init)
}

func (a *_apiIPTables) restore(c *gin.Context) {
	// var params data.Outbound
	// e := c.ShouldBindWith(&params, binding.JSON)
	// if e != nil {
	// 	c.String(http.StatusInternalServerError, e.Error())
	// 	return
	// }
	// var mSettings manipulator.Settings
	// iptables, e := mSettings.GetIPtables()
	// if e != nil {
	// 	c.String(http.StatusInternalServerError, e.Error())
	// 	return
	// }
	// if strings.TrimSpace(iptables.Clear) == "" {
	// 	c.String(http.StatusInternalServerError, "clear command nil")
	// 	return
	// }
	// renderCommand(c, iptables.Shell, iptables.Clear)
	// return
}
func (a *_apiIPTables) init(c *gin.Context) {
	// var params data.Outbound
	// e := c.ShouldBindWith(&params, binding.JSON)
	// if e != nil {
	// 	c.String(http.StatusInternalServerError, e.Error())
	// 	return
	// }
	// var mSettings manipulator.Settings
	// iptables, e := mSettings.GetIPtables()
	// if e != nil {
	// 	c.String(http.StatusInternalServerError, e.Error())
	// 	return
	// }
	// if strings.TrimSpace(iptables.Init) == "" {
	// 	c.String(http.StatusInternalServerError, `init command nil`)
	// 	return
	// }
	// if strings.TrimSpace(iptables.Clear) != "" {
	// 	var bufferError bytes.Buffer
	// 	var bufferOut bytes.Buffer
	// 	buffer := bytes.NewBufferString(iptables.Clear)
	// 	cmd := exec.Command(iptables.Shell)
	// 	cmd.Stdin = buffer
	// 	cmd.Stdout = &bufferOut
	// 	cmd.Stderr = &bufferError
	// 	e = cmd.Run()
	// 	if e != nil {
	// 		if bufferError.Len() != 0 {
	// 			e = errors.New(bufferError.String())
	// 		}
	// 		c.String(http.StatusInternalServerError, e.Error())
	// 		return
	// 	}
	// }
	// text, e := getTemplate("init", &params, iptables.Init)
	// if e != nil {
	// 	return
	// }
	// renderCommand(c, iptables.Shell, text)
	return
}
