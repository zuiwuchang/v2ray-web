package utils

import (
	"reflect"
	"unsafe"
)

// StringToBytes string to []byte
func StringToBytes(str string) (b []byte) {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	sliceHeader.Data = strHeader.Data
	sliceHeader.Len = strHeader.Len
	sliceHeader.Cap = strHeader.Len
	return b
}

// BytesToString []byte to string
func BytesToString(b []byte) (str string) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))

	strHeader.Data = sliceHeader.Data
	strHeader.Len = sliceHeader.Len
	return str
}
