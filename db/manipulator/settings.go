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
		e = bucket.Put(key, []byte(data.V2rayTemplate))
		if e != nil {
			return
		}
	}

	key = []byte(data.SettingsIPTables)
	val = bucket.Get(key)
	if val == nil {
		var tmp data.IPTables
		tmp.ResetDefault()
		val, e = tmp.Encoder()
		if e != nil {
			return
		}
		e = bucket.Put(key, val)
		if e != nil {
			return
		}
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

// GetIPtables 返回 iptables 設定
func (m Settings) GetIPtables() (iptables *data.IPTables, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SettingsBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SettingsBucket)
			return
		}
		val := bucket.Get([]byte(data.SettingsIPTables))
		var tmp data.IPTables
		if val == nil {
			tmp.ResetDefault()
			iptables = &tmp
		} else {
			e = tmp.Decode(val)
			if e != nil {
				return
			}
			iptables = &tmp
		}
		return
	})
	return
}

// PutIPtables 保存 iptables 設定
func (m Settings) PutIPtables(iptables *data.IPTables) (e error) {
	b, e := iptables.Encoder()
	if e != nil {
		return
	}
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SettingsBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SettingsBucket)
			return
		}
		e = bucket.Put([]byte(data.SettingsIPTables), b)
		return
	})
	return
}
