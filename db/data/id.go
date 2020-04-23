package data

import (
	"bytes"
	"encoding/gob"
)

// EncodeID .
func EncodeID(id uint64) (b []byte, e error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	e = encoder.Encode(id)
	if e == nil {
		b = buffer.Bytes()
	}
	return
}

// DecoderID .
func DecoderID(b []byte) (id uint64, e error) {
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	e = decoder.Decode(&id)
	return
}
