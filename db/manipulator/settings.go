package manipulator

import (
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/utils"
)

// Settings 設定
type Settings struct {
}

// Init 初始化 bucket
func (m Settings) Init(tx *bolt.Tx) (e error) {
	bucket, e := tx.CreateBucketIfNotExists([]byte(data.SettingsBucket))
	if e != nil {
		return
	}
	key := []byte(data.SettingsV2ray)
	val := bucket.Get(key)
	if val == nil {
		bucket.Put(key, []byte(data.V2rayTemplate))
	}
	return
}

// GetV2ray 返回 v2ray 設定
func (m Settings) GetV2ray() (text string, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SettingsBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SettingsBucket)
			return
		}
		val := bucket.Get([]byte(data.SettingsV2ray))
		if val != nil {
			text = utils.BytesToString(val)
		}
		return
	})
	return
}

// PutV2ray 保存 v2ray 設定
func (m Settings) PutV2ray(text string) (e error) {
	text = strings.TrimSpace(text)
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SettingsBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SettingsBucket)
			return
		}
		if text == "" {
			e = bucket.Delete([]byte(data.SettingsV2ray))
		} else {
			e = bucket.Put([]byte(data.SettingsV2ray), utils.StringToBytes(text))
		}
		return
	})
	return
}
