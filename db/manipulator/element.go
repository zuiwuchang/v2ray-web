package manipulator

import (
	"fmt"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
)

// Element 代理節點
type Element struct {
}

// Init 初始化 bucket
func (m Element) Init(tx *bolt.Tx) (e error) {
	_, e = tx.CreateBucketIfNotExists([]byte(data.ElementBucket))
	if e != nil {
		return
	}
	return
}

// List 返回 所有記錄
func (m Element) List() (result []*data.Element, subscription []*data.Subscription, e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		// 返回 組信息
		var mSubscription Subscription
		subscription, e = mSubscription.list(t)
		if e != nil {
			return
		}

		// 返回訂閱節點
		bucket := t.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.ElementBucket)
			return
		}
		e = bucket.ForEach(func(k, v []byte) error {
			bucket := bucket.Bucket(k)
			if bucket != nil {
				bucket.ForEach(func(k, v []byte) error {
					var node data.Element
					e := node.Decode(v)
					if e == nil {
						result = append(result, &node)
					} else {
						if ce := logger.Logger.Check(zap.WarnLevel, "Decode Element error"); ce != nil {
							ce.Write(
								zap.Error(e),
							)
						}
					}
					return nil
				})
			}
			return nil
		})
		return
	})
	return
}

// Puts 更新記錄
func (m Element) Puts(subscription uint64, outbounds []*data.Outbound) (result []data.Element, e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		// 返回 組信息
		var mSubscription Subscription
		_, e = mSubscription.get(t, subscription)
		if e != nil {
			return
		}
		bucket := t.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.ElementBucket)
			return
		}

		// 刪除組
		key, e := data.EncodeID(subscription)
		if e != nil {
			return
		}
		e = bucket.DeleteBucket(key)
		if e != nil && e != bolt.ErrBucketNotFound {
			return
		}

		// 創建新組
		bucket, e = bucket.CreateBucket(key)
		if e != nil {
			return
		}
		// 插入記錄
		count := len(outbounds)
		arrs := make([]data.Element, count)
		var val []byte
		for i := 0; i < count; i++ {
			arrs[i].ID, e = bucket.NextSequence()
			if e != nil {
				return
			}
			key, e = data.EncodeID(arrs[i].ID)
			if e != nil {
				return
			}
			arrs[i].Outbound = outbounds[i]
			arrs[i].Subscription = subscription
			val, e = arrs[i].Encoder()
			if e != nil {
				return
			}

			e = bucket.Put(key, val)
			if e != nil {
				return
			}
		}
		result = arrs
		return
	})
	return
}
