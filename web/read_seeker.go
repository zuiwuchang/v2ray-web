package web

import (
	"bytes"
	"io"
)

// ReadSeeker .
type ReadSeeker struct {
	reader  *bytes.Reader
	Data    interface{}
	Marshal func(in interface{}) (out []byte, err error)
}

// Read .
func (r *ReadSeeker) Read(p []byte) (n int, err error) {
	err = r.Check()
	if err != nil {
		return
	}
	n, err = r.reader.Read(p)
	return
}

// WriteTo .
func (r *ReadSeeker) WriteTo(w io.Writer) (n int64, err error) {
	err = r.Check()
	if err != nil {
		return
	}
	n, err = r.reader.WriteTo(w)
	return
}

// Seek .
func (r *ReadSeeker) Seek(offset int64, whence int) (abs int64, err error) {
	err = r.Check()
	if err != nil {
		return 0, err
	}
	abs, err = r.reader.Seek(offset, whence)
	return
}

// Check .
func (r *ReadSeeker) Check() (e error) {
	if r.reader != nil {
		return
	}
	b, e := r.Marshal(r.Data)
	if e != nil {
		return
	}
	r.reader = bytes.NewReader(b)
	return
}
