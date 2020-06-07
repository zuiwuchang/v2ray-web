package configure

import (
	"path/filepath"
	"strings"
	"time"
)

// HTTP configure http
type HTTP struct {
	Addr string

	CertFile string
	KeyFile  string

	// cookie 過期時間
	Maxage time.Duration
	// cookie 密鑰
	Secret string
}

// Safe if tls return true
func (c *HTTP) Safe() bool {
	return c.CertFile != "" && c.KeyFile != ""
}

// Format .
func (c *HTTP) Format(basePath string) (e error) {
	c.Addr = strings.TrimSpace(c.Addr)
	c.CertFile = strings.TrimSpace(c.CertFile)
	c.KeyFile = strings.TrimSpace(c.KeyFile)
	c.Secret = strings.TrimSpace(c.Secret)

	if c.Safe() {
		if filepath.IsAbs(c.CertFile) {
			c.CertFile = filepath.Clean(c.CertFile)
		} else {
			c.CertFile = filepath.Clean(basePath + "/" + c.CertFile)
		}

		if filepath.IsAbs(c.KeyFile) {
			c.KeyFile = filepath.Clean(c.KeyFile)
		} else {
			c.KeyFile = filepath.Clean(basePath + "/" + c.KeyFile)
		}
	}
	if c.Maxage > 0 {
		c.Maxage *= time.Millisecond
	} else {
		c.Maxage = time.Hour * 24 * 30
	}

	return
}
