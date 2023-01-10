package manipulator

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/db/data"
)

// Version 數據庫 當前版本
const Version = 7

// manipulator 定義數據庫 操縱器 接口
type manipulator interface {
	// 初始化數據庫
	Init(tx *bolt.Tx, version int) (e error)
	// 升級數據庫
	Upgrade(tx *bolt.Tx, oldVersion, newVersion int) (e error)
}

func updateVersion(tx *bolt.Tx) (oldVersion int, e error) {
	bucketName := []byte(`__private_data`)
	bucket, e := tx.CreateBucketIfNotExists(bucketName)
	if e != nil {
		return
	}

	// 獲取舊版本
	keyVersion := []byte(`version`)
	b := bucket.Get(keyVersion)
	if len(b) == 4 {
		oldVersion = int(binary.LittleEndian.Uint32(b))
	}
	if oldVersion == 0 {
		bucket := tx.Bucket([]byte(data.UserBucket))
		if bucket != nil {
			oldVersion = 1
		}
	}
	// 設置新版本
	if oldVersion > Version {
		e = errors.New(`the local database version is greater than the current version`)
		return
	} else if Version > oldVersion {
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, Version)
		e = bucket.Put(keyVersion, b)
	}
	return
}
