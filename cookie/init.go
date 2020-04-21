package cookie

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gitlab.com/king011/v2ray-web/utils"

	"gitlab.com/king011/v2ray-web/logger"

	"github.com/gorilla/securecookie"
	"go.uber.org/zap"
)

var _Secure *securecookie.SecureCookie

// MaxAge return max age of timeout
func MaxAge() int64 {
	return _MaxAge
}

// Encode encode cookie
func Encode(name string, value interface{}) (string, error) {
	return _Secure.Encode(name, value)
}

// Decode decode  cookie
func Decode(name, value string, dst interface{}) error {
	return _Secure.Decode(name, value, dst)
}

// IsInit return is init ?
func IsInit() bool {
	return _Secure != nil
}

var _MaxAge int64

// Init initialize cookie system
func Init(secret string, maxAge time.Duration) (e error) {
	b := md5.Sum(utils.StringToBytes(secret))
	var hashKey [32]byte
	hex.Encode(hashKey[:], b[:])
	blockKey := sha256.Sum256(utils.StringToBytes(secret))
	initKey(hashKey[:], blockKey[:], maxAge)
	return
}
func initKey(hashKey, blockKey []byte, maxAge time.Duration) {
	_Secure = securecookie.New(hashKey, blockKey)
	_MaxAge = int64(maxAge / time.Second)
	_Secure.MaxAge((int)(_MaxAge))
	if ce := logger.Logger.Check(zap.InfoLevel, "cookie"); ce != nil {
		ce.Write(
			zap.String("timeout", maxAge.String()),
		)
	}
	return
}
