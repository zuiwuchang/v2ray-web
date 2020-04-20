package logger

import (
	"log"
	"path/filepath"
	"strings"

	logger "gitlab.com/king011/king-go/log/logger.zap"
	"go.uber.org/zap"
)

// Logger logger single
var Logger logger.Logger

// JoinFields .
var JoinFields = logger.JoinFields

// Fields .
var Fields = logger.Fields

// Init init logger
func Init(basePath string, options *logger.Options) (e error) {
	// format
	options.Filename = strings.TrimSpace(options.Filename)
	if options.Filename != "" {
		if !filepath.IsAbs(options.Filename) {
			options.Filename = filepath.Clean(basePath + "/" + options.Filename)
		}
	}
	var zapOptions []zap.Option
	if options.Caller {
		zapOptions = append(zapOptions, zap.AddCaller())
	}

	// new logger
	l := logger.New(options, zapOptions...)

	// run http
	if options.HTTP != "" {
		errHTTP := l.StartHTTP()
		if errHTTP == nil {
			if l.OutFile() {
				log.Println("zap http running", options.HTTP)
			}
			l.Info("zap http",
				zap.String("running", options.HTTP),
			)
		} else {
			if l.OutFile() {
				log.Println("zap http running", errHTTP)
			}
			l.Warn("zap http",
				zap.Error(errHTTP),
			)
		}
	}
	// Attach
	Logger.Attach(l)
	return
}
