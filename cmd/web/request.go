package web

import (
	"encoding/json"
	"net/http"

	"gitlab.com/king011/v2ray-web/utils"
)

// Helper http 輔助類
type Helper struct {
	request  *http.Request
	body     []byte
	response http.ResponseWriter
}

// BodyJSON 以 json 解碼 request body
func (h *Helper) BodyJSON(v interface{}) (e error) {
	e = json.Unmarshal(h.body, v)
	return e
}

// RenderJSON 以 json 響應
func (h *Helper) RenderJSON(v interface{}) {
	header := h.response.Header()
	header.Set("Content-Type", "application/json; charset=UTF-8")

	b, e := json.Marshal(v)
	if e != nil {
		return
	}
	h.response.WriteHeader(http.StatusOK)
	h.response.Write(b)
	return
}

// RenderError 以 錯誤 響應
func (h *Helper) RenderError(e error) {
	if e == nil {
		h.RenderText(http.StatusInternalServerError, "")
	} else {
		h.RenderText(http.StatusInternalServerError, e.Error())
	}
}

// RenderText 以 文本 響應
func (h *Helper) RenderText(statusCode int, text string) {
	if text != "" {
		header := h.response.Header()
		header.Set("Content-Type", "text/plain; charset=UTF-8")
	}
	h.response.WriteHeader(statusCode)
	if text != "" {
		h.response.Write(utils.StringToBytes(text))
	}
}
