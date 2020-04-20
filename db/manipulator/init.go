package manipulator

import (
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

	_db = db
	return
}
