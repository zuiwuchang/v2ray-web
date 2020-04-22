package manipulator

import (
	"os"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/configure"
	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
)

type manipulator interface {
	Init(tx *bolt.Tx) (e error)
}

var _db *bolt.DB

// Init 初始化 數據庫
func Init(cnf *configure.Database) (e error) {
	db, e := bolt.Open(cnf.Source, 0600, nil)
	if e != nil {
		if ce := logger.Logger.Check(zap.FatalLevel, "open databases error"); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String("source", cnf.Source),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, "open databases"); ce != nil {
		ce.Write(
			zap.String("source", cnf.Source),
		)
	}

	e = db.Update(func(tx *bolt.Tx) (e error) {
		buckets := []manipulator{
			User{},
			Settings{},
			Subscription{},
		}
		for i := 0; i < len(buckets); i++ {
			e = buckets[i].Init(tx)
			if e != nil {
				return
			}
		}
		return
	})
	if e != nil {
		os.Exit(1)
	}
	_db = db
	return
}

// DB .
func DB() *bolt.DB {
	return _db
}
