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
func (m Element) Init(tx *bolt.Tx, version int) (e error) {
	bucket, e := tx.CreateBucketIfNotExists([]byte(data.ElementBucket))
	if e != nil {
		return
	}
	key, e := data.EncodeID(0)
	if e != nil {
		return
	}
	_, e = bucket.CreateBucketIfNotExists(key)
	return
}

// Upgrade 升級 bucket
func (m Element) Upgrade(tx *bolt.Tx, oldVersion, newVersion int) (e error) {
	e = m.Init(tx, newVersion)
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

// Add 增加記錄
func (Element) Add(subscription uint64, outbound *data.Outbound) (result uint64, e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.ElementBucket)
			return
		}

		key, e := data.EncodeID(subscription)
		if e != nil {
			return
		}
		bucket = bucket.Bucket(key)
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s.%v", data.ElementBucket, subscription)
			return
		}
		id, e := bucket.NextSequence()
		if e != nil {
			return
		}
		// 插入記錄
		key, e = data.EncodeID(id)
		if e != nil {
			return
		}
		element := data.Element{
			ID:           id,
			Subscription: subscription,
			Outbound:     outbound,
		}
		val, e := element.Encoder()
		if e != nil {
			return
		}
		e = bucket.Put(key, val)
		if e != nil {
			return
		}
		result = id
		return
	})
	return
}

// Put 更新記錄
func (Element) Put(subscription, id uint64, outbound *data.Outbound) (e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.ElementBucket)
			return
		}

		key, e := data.EncodeID(subscription)
		if e != nil {
			return
		}
		bucket = bucket.Bucket(key)
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s.%v", data.ElementBucket, subscription)
			return
		}
		// 查找記錄
		key, e = data.EncodeID(id)
		if e != nil {
			return
		}
		val := bucket.Get(key)
		if val == nil {
			e = fmt.Errorf("key not exist : %s.%v%v", data.ElementBucket, subscription, id)
			return
		}
		element := data.Element{
			ID:           id,
			Subscription: subscription,
			Outbound:     outbound,
		}
		val, e = element.Encoder()
		if e != nil {
			return
		}
		e = bucket.Put(key, val)
		if e != nil {
			return
		}
		return
	})
	return
}

// Remove 刪除記錄
func (Element) Remove(subscription, id uint64) (e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.ElementBucket)
			return
		}

		key, e := data.EncodeID(subscription)
		if e != nil {
			return
		}
		bucket = bucket.Bucket(key)
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s.%v", data.ElementBucket, subscription)
			return
		}
		// 查找記錄
		key, e = data.EncodeID(id)
		if e != nil {
			return
		}
		val := bucket.Get(key)
		if val == nil {
			e = fmt.Errorf("key not exist : %s.%v%v", data.ElementBucket, subscription, id)
			return
		}
		e = bucket.Delete(key)
		if e != nil {
			return
		}
		return
	})
	return
}

// Clear 清空記錄
func (Element) Clear(subscription uint64) (e error) {
	e = _db.Update(func(t *bolt.Tx) (e error) {
		bucket := t.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.ElementBucket)
			return
		}

		key, e := data.EncodeID(subscription)
		if e != nil {
			return
		}
		e = bucket.DeleteBucket(key)
		if e != nil {
			if e == bolt.ErrBucketNotFound {
				e = fmt.Errorf("bucket not exist : %s.%v", data.ElementBucket, subscription)
			}
			return
		}
		_, e = bucket.CreateBucket(key)
		if e != nil {
			return
		}
		return
	})
	return
}
func (m Element) Import(vals []*data.Element) error {
	return _db.Update(func(tx *bolt.Tx) (e error) {
		bucket := tx.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.StrategyBucket)
			return
		}
		c := bucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if e = bucket.DeleteBucket(k); e != nil {
				return e
			}
		}
		if len(vals) == 0 {
			return
		}
		im := newImportElement(tx)
		for _, val := range vals {
			e = im.put(val)
			if e != nil {
				return
			}
		}
		return
	})
}

type importElement struct {
	tx     *bolt.Tx
	bucket *bolt.Bucket
	keys   map[uint64]*bolt.Bucket
}

func newImportElement(tx *bolt.Tx) *importElement {
	return &importElement{
		tx:   tx,
		keys: make(map[uint64]*bolt.Bucket),
	}
}

func (i *importElement) getBucket(key uint64) (bucket *bolt.Bucket, e error) {
	bucket, ok := i.keys[key]
	if ok {
		return
	}
	bucket = i.bucket
	if bucket == nil {
		bucket = i.tx.Bucket([]byte(data.ElementBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.ElementBucket)
			return
		}
		i.bucket = bucket
	}

	b, e := data.EncodeID(key)
	if e != nil {
		return
	}
	bucket, e = bucket.CreateBucketIfNotExists(b)
	if e != nil {
		return
	}
	i.keys[key] = bucket
	i.bucket = bucket
	return
}

func (i *importElement) put(ele *data.Element) (e error) {
	bucket, e := i.getBucket(ele.Subscription)
	if e != nil {
		return
	}
	key, e := data.EncodeID(ele.ID)
	if e != nil {
		return
	}
	val, e := ele.Encoder()
	if e != nil {
		return
	}
	bucket.Put(key, val)
	return
}
