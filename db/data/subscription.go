package data

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(Subscription{})
}

// SubscriptionBucket .
const SubscriptionBucket = "subscription"

// Subscription 訂閱服務
type Subscription struct {
	// 唯一識別碼
	ID uint64 `json:"id,omitempty"`
	// 給人類看的名稱
	Name string `json:"name,omitempty"`
	// 訂閱地址
	URL string `json:"url,omitempty"`
}

// Decode 由 []byte 解碼
func (s *Subscription) Decode(b []byte) (e error) {
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	e = decoder.Decode(s)
	return
}

// Encoder 編碼到 []byte
func (s *Subscription) Encoder() (b []byte, e error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	e = encoder.Encode(s)
	if e != nil {
		b = buffer.Bytes()
	}
	return
}
