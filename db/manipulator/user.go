package manipulator

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
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

// List 返回用戶 列表
func (m User) List() (result []data.User, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket(data.UserBucket)
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.UserBucket)
			return
		}
		bucket.ForEach(func(k, v []byte) error {
			result = append(result, data.User{
				Name:     utils.BytesToString(k),
				Password: utils.BytesToString(v),
			})
			return nil
		})
		return
	})
	return
}

// Add 添加用戶
func (m User) Add(name, password string) (e error) {
	if name == "" {
		e = errors.New("name not support empty")
		return
	}
	if password == "" {
		e = errors.New("password not support empty")
		return
	}
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket(data.UserBucket)
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.UserBucket)
			return
		}
		key := utils.StringToBytes(name)
		v := bucket.Get(key)
		if v != nil {
			e = fmt.Errorf("user already exists : %s", name)
			return
		}
		e = bucket.Put(key, utils.StringToBytes(password))
		return
	})
	return
}

// Remove 刪除用戶
func (m User) Remove(name string) (e error) {
	if name == "" {
		e = errors.New("name not support empty")
		return
	}
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket(data.UserBucket)
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.UserBucket)
			return
		}
		key := utils.StringToBytes(name)
		e = bucket.Delete(key)
		return
	})
	return
}

// Password 修改密碼
func (m User) Password(name, password string) (e error) {
	if name == "" {
		e = errors.New("name not support empty")
		return
	}
	if password == "" {
		e = errors.New("password not support empty")
		return
	}
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket(data.UserBucket)
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.UserBucket)
			return
		}
		key := utils.StringToBytes(name)
		v := bucket.Get(key)
		if v == nil {
			e = fmt.Errorf("user not exists : %s", name)
			return
		}
		e = bucket.Put(key, utils.StringToBytes(password))
		return
	})
	return
}
