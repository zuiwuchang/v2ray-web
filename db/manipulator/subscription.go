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
func (m Subscription) Init(tx *bolt.Tx, version int) (e error) {
	_, e = tx.CreateBucketIfNotExists([]byte(data.SubscriptionBucket))
	if e != nil {
		return
	}
	return
}

// Upgrade 升級 bucket
func (m Subscription) Upgrade(tx *bolt.Tx, oldVersion, newVersion int) (e error) {
	e = m.Init(tx, newVersion)
	return
}
func (m Subscription) list(t *bolt.Tx) (result []*data.Subscription, e error) {
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
}

// List 返回 所有記錄
func (m Subscription) List() (result []*data.Subscription, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		result, e = m.list(t)
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

// Get 返回記錄
func (m Subscription) Get(id uint64) (result *data.Subscription, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		result, e = m.get(t, id)
		return
	})
	return
}
func (m Subscription) get(t *bolt.Tx, id uint64) (result *data.Subscription, e error) {
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
		return
	}
	result = &node
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

// Remove 刪除記錄
func (m Subscription) Remove(id uint64) (e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.SubscriptionBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.SubscriptionBucket)
			return
		}
		key, e := data.EncodeID(id)
		if e != nil {
			return
		}
		e = bucket.Delete(key)
		if e != nil {
			return
		}
		// 刪除 訂閱
		bucket = t.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.ElementBucket)
			return
		}
		e = bucket.DeleteBucket(key)
		if e != nil {
			if e == bolt.ErrBucketNotFound {
				e = nil
			}
			return
		}
		return
	})
	return
}
func (m Subscription) Import(vals []*data.Subscription) error {
	return _db.Update(func(tx *bolt.Tx) (e error) {
		bucket := tx.Bucket([]byte(data.SubscriptionBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.StrategyBucket)
			return
		}
		c := bucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if e = bucket.Delete(k); e != nil {
				return e
			}
		}

		var b, key []byte
		for _, val := range vals {
			b, e = val.Encoder()
			if e != nil {
				return
			}
			key, e = data.EncodeID(val.ID)
			if e != nil {
				return
			}
			e = bucket.Put(key, b)
			if e != nil {
				return
			}
		}
		return
	})
}
