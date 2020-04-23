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

// List 返回 所有記錄
func (m Element) List() (result []*data.Element, subscription []*data.Subscription, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
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
		return
	})
	return
}
