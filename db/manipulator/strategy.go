package manipulator

import (
	"fmt"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/utils"
	"go.uber.org/zap"
)

type Strategy struct {
}

func (m Strategy) init(tx *bolt.Tx) (bucket *bolt.Bucket, e error) {
	bucket, e = tx.CreateBucketIfNotExists([]byte(data.StrategyBucket))
	if e != nil {
		return
	}
	key := []byte(data.StrategyDefault)
	val := bucket.Get(key)
	if len(val) != 0 {
		var s data.Strategy
		if s.Decode(val) == nil {
			return
		}
	}

	items := []utils.Pair[string, int]{
		utils.MakePair(data.StrategyDefault, 0),
		utils.MakePair(`All Proxy`, 1),
		utils.MakePair(`Public IP Proxy`, 100),
		utils.MakePair(`Proxy First`, 200),
		utils.MakePair(`Direct First`, 900),
		utils.MakePair(`All Direct`, 1000),
	}
	for _, item := range items {
		s := data.Strategy{
			Name:  item.First,
			Value: item.Second,
		}
		val, err := s.Encoder()
		if err != nil {
			e = err
			return
		}
		e = bucket.Put([]byte(s.Name), val)
	}
	return
}
func (m Strategy) Init(tx *bolt.Tx, version int) (e error) {
	_, e = m.init(tx)
	if e != nil {
		return
	}
	return
}
func (m Strategy) Upgrade(tx *bolt.Tx, oldVersion, newVersion int) (e error) {
	_, e = m.init(tx)
	if e != nil {
		return
	}
	return
}
func (m Strategy) List() (result []*data.Strategy, e error) {
	e = _db.View(func(t *bolt.Tx) (e error) {
		result, e = m.list(t)
		return
	})
	return
}
func (m Strategy) list(tx *bolt.Tx) (result []*data.Strategy, e error) {
	bucket := tx.Bucket([]byte(data.StrategyBucket))
	if bucket == nil {
		e = fmt.Errorf("bucket not exist : %s", data.StrategyBucket)
		return
	}
	e = bucket.ForEach(func(k, v []byte) error {
		var node data.Strategy
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

func (m Strategy) Value(name string) (result *data.StrategyValue, e error) {
	e = _db.View(func(tx *bolt.Tx) (e error) {
		result, e = m.value(tx, name)
		return
	})
	return
}
func (m Strategy) value(tx *bolt.Tx, name string) (result *data.StrategyValue, e error) {
	bucket := tx.Bucket([]byte(data.StrategyBucket))
	if bucket == nil {
		e = fmt.Errorf("bucket not exist : %s", data.StrategyBucket)
		return
	}
	def, e := m.getValue(bucket, []byte(data.StrategyDefault))
	if e != nil {
		return
	}
	if name == `` || name == data.StrategyDefault {
		def.Host = m.copyHost(def.Host, nil)
		def.ProxyIP = m.copy(def.ProxyIP, nil)
		def.ProxyDomain = m.copy(def.ProxyDomain, nil)
		def.DirectIP = m.copy(def.DirectIP, nil)
		def.DirectDomain = m.copy(def.DirectDomain, nil)
		def.BlockIP = m.copy(def.BlockIP, nil)
		def.BlockDomain = m.copy(def.BlockDomain, nil)

		result = def
		return
	}
	val, e := m.getValue(bucket, []byte(name))
	if e != nil {
		return
	}
	val.Host = m.copyHost(val.Host, def.Host)
	val.ProxyIP = m.copy(val.ProxyIP, def.ProxyIP)
	val.ProxyDomain = m.copy(val.ProxyDomain, def.ProxyDomain)
	val.DirectIP = m.copy(val.DirectIP, def.DirectIP)
	val.DirectDomain = m.copy(val.DirectDomain, def.DirectDomain)
	val.BlockIP = m.copy(val.BlockIP, def.BlockIP)
	val.BlockDomain = m.copy(val.BlockDomain, def.BlockDomain)

	result = val
	return
}
func (m Strategy) copyHost(dst [][]string, src [][]string) [][]string {
	out := make([][]string, 0, len(dst)+len(src))
	keys := make(map[string]bool)
	for i := len(dst) - 1; i >= 0; i-- {
		h := dst[i]
		if len(h) < 2 || h[0] == `` || keys[h[0]] {
			continue
		}

		keys[h[0]] = true
		out = append(out, h)
	}

	for i := len(src) - 1; i >= 0; i-- {
		h := src[i]
		if len(h) < 2 || h[0] == `` || keys[h[0]] {
			continue
		}
		keys[h[0]] = true
		out = append(out, h)
	}
	return out
}
func (m Strategy) copy(dst []string, src []string) []string {
	out := make([]string, 0, len(dst)+len(src))
	keys := make(map[string]bool)
	for _, key := range src {
		if key == `` || keys[key] {
			continue
		}
		keys[key] = true
		out = append(out, key)
	}
	for _, key := range dst {
		if key == `` || keys[key] {
			continue
		}
		keys[key] = true
		out = append(out, key)
	}
	return out
}
func (m Strategy) get(bucket *bolt.Bucket, key []byte) (result *data.Strategy, e error) {
	b := bucket.Get(key)
	if len(b) == 0 {
		if string(key) == data.StrategyDefault {
			result = &data.Strategy{
				Name: data.StrategyDefault,
			}
			return
		}
		e = fmt.Errorf("key not exist : %s.%v", data.StrategyBucket, key)
		return
	}
	var s data.Strategy
	e = s.Decode(b)
	if e != nil {
		return
	}
	result = &s
	return
}
func (m Strategy) getValue(bucket *bolt.Bucket, key []byte) (result *data.StrategyValue, e error) {
	v, e := m.get(bucket, key)
	if e != nil {
		return
	}
	result = v.ToValue()
	return
}
func (m Strategy) Put(d *data.Strategy) error {
	if d.Name == data.StrategyDefault && d.Value != 0 {
		return fmt.Errorf("'%s.value' must be 0", data.StrategyDefault)
	}

	value, e := d.Encoder()
	if e != nil {
		return e
	}
	return _db.Update(func(tx *bolt.Tx) (e error) {
		bucket := tx.Bucket([]byte(data.StrategyBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.StrategyBucket)
			return
		}
		b := bucket.Get([]byte(d.Name))
		if len(b) == 0 {
			e = fmt.Errorf("strategy not exist : %s", d.Name)
			return
		}
		e = bucket.Put([]byte(d.Name), value)
		return
	})
}
func (m Strategy) Add(d *data.Strategy) error {
	if d.Name == data.StrategyDefault && d.Value != 0 {
		return fmt.Errorf("'%s.value' must be 0", data.StrategyDefault)
	}

	value, e := d.Encoder()
	if e != nil {
		return e
	}
	return _db.Update(func(tx *bolt.Tx) (e error) {
		bucket := tx.Bucket([]byte(data.StrategyBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.StrategyBucket)
			return
		}
		b := bucket.Get([]byte(d.Name))
		if len(b) != 0 {
			e = fmt.Errorf("strategy already exist : %s", d.Name)
			return
		}
		e = bucket.Put([]byte(d.Name), value)
		return
	})
}
func (m Strategy) Remove(name string) error {
	if name == data.StrategyDefault {
		return fmt.Errorf("'%s' cannot be deleted", name)
	}

	return _db.Update(func(tx *bolt.Tx) (e error) {
		bucket := tx.Bucket([]byte(data.StrategyBucket))
		if bucket == nil {
			e = fmt.Errorf("bucket not exist : %s", data.StrategyBucket)
			return
		}
		e = bucket.Delete([]byte(name))
		return
	})
}
