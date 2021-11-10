package logger

import (
	logger "gitlab.com/king011/king-go/log/logger.zap"
	"go.uber.org/zap"
)

func InitConsole(level string) {
	lv := zap.NewAtomicLevel()
	if level == `` {
		lv.SetLevel(zap.InfoLevel)
	} else if e := lv.UnmarshalText([]byte(level)); e != nil {
		lv.SetLevel(zap.InfoLevel)
	}
	Logger.Attach(logger.New(&logger.Options{
		Level:  lv.String(),
		Caller: true,
	}))
}
