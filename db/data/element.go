package data

import (
	"bytes"
	"encoding/gob"
)

// ElementBucket .
const ElementBucket = "element"

func init() {
	gob.Register(Element{})
}

// Element 代理節點
type Element struct {
	// 唯一識別碼
	ID uint64 `json:"id,omitempty"`
	// 所屬的 訂閱組 如果爲0 則爲 自定義節點
	Subscription uint64 `json:"subscription,omitempty"`
	// 節點信息
	Outbound Outbound `json:"outbound,omitempty"`
}

// Decode 由 []byte 解碼
func (element *Element) Decode(b []byte) (e error) {
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	e = decoder.Decode(element)
	return
}

// Encoder 編碼到 []byte
func (element *Element) Encoder() (b []byte, e error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	e = encoder.Encode(element)
	if e == nil {
		b = buffer.Bytes()
	}
	return
}
