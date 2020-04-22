package manipulator

import (
	"fmt"

	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/db/data"
)

// Subscription .
type Subscription struct {
}

// Init 初始化 bucket
func (m Subscription) Init(tx *bolt.Tx) (e error) {
	_, e = tx.CreateBucketIfNotExists([]byte(data.SubscriptionBucket))
	if e != nil {
		return
	}
	return
}

// List 返回 所有記錄
func (m Subscription) List() (result []*data.Subscription, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SubscriptionBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SubscriptionBucket)
			return
		}
		e = bucket.ForEach(func(k, v []byte) error {
			var node data.Subscription
			e := node.Decode(v)
			if e == nil {
				result = append(result, &node)
			} else {
				if ce := logger.Logger.Check(zap.WarnLevel, "Decode Subscription error"); ce != nil {
					ce.Write(
						zap.Error(e),
					)
				}
			}
			return nil
		})
		return
	})
	return
}

// Get 返回記錄
func (m Subscription) Get(id uint64) (result *data.Subscription, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SubscriptionBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SubscriptionBucket)
			return
		}
		key, e := data.EncodeID(id)
		if e != nil {
			return
		}
		val := bucket.Get(key)
		if val == nil {
			e = fmt.Errorf("key not exist : %s.%v", data.SubscriptionBucket, id)
			return
		}
		var node data.Subscription
		e = node.Decode(val)
		if e != nil {
			result = &node
		}
		return
	})
	return
}

// Put 設置記錄
func (m Subscription) Put(node *data.Subscription) (e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SubscriptionBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SubscriptionBucket)
			return
		}
		key, e := data.EncodeID(node.ID)
		if e != nil {
			return
		}
		val := bucket.Get(key)
		if val == nil {
			e = fmt.Errorf("key not exist : %s.%v", data.SubscriptionBucket, node.ID)
			return
		}
		val, e = node.Encoder()
		if e != nil {
			return
		}
		e = bucket.Put(key, val)
		return
	})
	return
}

// Add 添加記錄
func (m Subscription) Add(node *data.Subscription) (e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SubscriptionBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SubscriptionBucket)
			return
		}
		id, e := bucket.NextSequence()
		if e != nil {
			return
		}
		key, e := data.EncodeID(id)
		if e != nil {
			return
		}
		node.ID = id
		val, e := node.Encoder()
		if e != nil {
			return
		}
		e = bucket.Put(key, val)
		return
	})
	return
}
