package manipulator

import (
	"fmt"

	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/utils"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/db/data"
)

// User 用戶 操縱器
type User struct {
}

// Login 登入
func (m User) Login(name, password string) (result *cookie.Session, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket(data.UserBucket)
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.UserBucket)
			return
		}
		val := bucket.Get(utils.StringToBytes(name))
		if val == nil {
			return
		}
		if utils.BytesToString(val) == password {
			result = &cookie.Session{
				Name: name,
				Root: true,
			}
		}
		return
	})
	return
}
