package manipulator

import (
	"os"

	"github.com/boltdb/bolt"
	"gitlab.com/king011/v2ray-web/configure"
	"gitlab.com/king011/v2ray-web/logger"
	"go.uber.org/zap"
)

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
		var oldVersion int
		oldVersion, e = updateVersion(tx)
		if e != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "database is not compatible"); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
			return
		}
		if oldVersion < 0 {
			oldVersion = 0
		}
		buckets := []manipulator{
			User{},
			Settings{},
			Subscription{},
			Element{},
			Strategy{},
		}
		if oldVersion == 0 {
			for i := 0; i < len(buckets); i++ {
				e = buckets[i].Init(tx, Version)
				if e != nil {
					return
				}
			}
		} else if oldVersion < Version {
			for i := 0; i < len(buckets); i++ {
				e = buckets[i].Upgrade(tx, oldVersion, Version)
				if e != nil {
					return
				}
			}
		}
		var m Strategy
		e = m.Upgrade(tx, oldVersion, Version)
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
