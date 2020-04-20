package configure

import (
	"path/filepath"
	"strings"
)

// Database 數據庫配置
type Database struct {
	// 數據源 位置
	Source string
}

// Format .
func (d *Database) Format(basePath string) (e error) {
	d.Source = strings.TrimSpace(d.Source)
	if d.Source == "" {
		d.Source = "v2ray-web.db"
	}
	if filepath.IsAbs(d.Source) {
		d.Source = filepath.Clean(d.Source)
	} else {
		d.Source = filepath.Clean(basePath + "/" + d.Source)
	}
	return
}
