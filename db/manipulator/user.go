package manipulator

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/utils"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/db/data"
)

// User 用戶 操縱器
type User struct {
}

// Init 初始化 bucket
func (m User) Init(tx *bolt.Tx) (e error) {
	bucket, e := tx.CreateBucketIfNotExists(data.UserBucket)
	if e != nil {
		return
	}
	cursor := bucket.Cursor()
	k, _ := cursor.First()
	if k != nil {
		return
	}
	password := sha512.Sum512([]byte("19890604"))
	e = bucket.Put([]byte("killer"), utils.StringToBytes(hex.EncodeToString(password[:])))
	return
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
