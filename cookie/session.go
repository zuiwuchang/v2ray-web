package cookie

import (
	"encoding/gob"
)

const (
	// CookieName cookie key name
	CookieName = "v2ray_web_session"
)

func init() {
	gob.Register(&Session{})
}

// Session user session info
type Session struct {
	Name string
	Root bool
}

// Cookie encode to cookie
func (s *Session) Cookie() (string, error) {
	return Encode("session", s)
}

// FromCookie restore session from cookie
func FromCookie(val string) (session *Session, e error) {
	var s Session
	e = Decode("session", val, &s)
	if e != nil {
		return
	}
	session = &s
	return
}
