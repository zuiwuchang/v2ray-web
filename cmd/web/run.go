package web

import (
	"bytes"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gitlab.com/king011/v2ray-web/configure"
	"gitlab.com/king011/v2ray-web/internal/logs"
	"gitlab.com/king011/v2ray-web/logger"
)

// Run .
func Run(cnf *configure.Configure, debug bool) {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	l, e := net.Listen("tcp", cnf.HTTP.Addr)
	if e == nil {
		if ce := logger.Logger.Check(zap.InfoLevel, "listen success"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.HTTP.Addr),
			)
		}
	} else {
		if ce := logger.Logger.Check(zap.InfoLevel, "listen error"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.HTTP.Addr),
			)
		}
		os.Exit(1)
	}
	defer l.Close()
	server, e := NewServer(l)
	if e != nil {
		os.Exit(1)
	}
	stdout, e := pipeBuffer(os.Stdout)
	if e == nil {
		os.Stdout = stdout.pw
		go func() {
			stdout.run()
			os.Stdout = stdout.source
		}()
	} else {
		if ce := logger.Logger.Check(zap.WarnLevel, "stdout pipe error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
	}
	stderr, e := pipeBuffer(os.Stderr)
	if e == nil {
		os.Stderr = stderr.pw
		go func() {
			stderr.run()
			os.Stderr = stderr.source
		}()
	} else {
		if ce := logger.Logger.Check(zap.WarnLevel, "stderr pipe error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
	}
	server.onStart()

	if cnf.HTTP.Safe() {
		if ce := logger.Logger.Check(zap.InfoLevel, "https serve"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.HTTP.Addr),
			)
		}
		e = server.ServeTLS(cnf.HTTP.CertFile, cnf.HTTP.KeyFile)
	} else {
		if ce := logger.Logger.Check(zap.InfoLevel, "http serve"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.HTTP.Addr),
			)
		}
		e = server.Serve()
	}
	if e != nil {
		if ce := logger.Logger.Check(zap.FatalLevel, "serve error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
	}
}

// Buffer .
type _Buffer struct {
	source *os.File
	pr     *os.File
	pw     *os.File
	leave  string
}

func pipeBuffer(source *os.File) (buffer *_Buffer, e error) {
	pr, pw, e := os.Pipe()
	if e != nil {
		return
	}
	buffer = &_Buffer{
		source: source,
		pr:     pr,
		pw:     pw,
	}
	return
}
func (b *_Buffer) run() {
	buffer := make([]byte, 4096)
	for {
		n, e := b.pr.Read(buffer)
		if e != nil {
			break
		} else if n == 0 {
			continue
		}
		b.source.Write(buffer[:n])
		b.push(buffer[:n])
	}
}
func (b *_Buffer) push(buffer []byte) {
	count := len(buffer)
	if buffer[count-1] == '\n' {
		bs := bytes.Split(buffer, []byte("\n"))
		count = len(bs)
		for i := 0; i < count; i++ {
			str := string(bs[i])
			if i == 0 && b.leave != "" {
				str = b.leave + str
				b.leave = ""
			}
			logs.Push(str)
		}
	} else {
		bs := bytes.Split(buffer, []byte("\n"))
		count = len(bs)
		for i := 0; i < count-1; i++ {
			str := string(bs[i])
			if i == 0 && b.leave != "" {
				str = b.leave + str
				b.leave = ""
			}
			logs.Push(str)
		}
		b.leave += string(bs[count-1])
		if len(b.leave) > 2048 {
			logs.Push(b.leave)
			b.leave = ""
		}
	}
}
