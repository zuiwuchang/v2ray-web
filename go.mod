module gitlab.com/king011/v2ray-web

go 1.12

require (
	github.com/boltdb/bolt v1.3.1
	github.com/gin-contrib/gzip v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/google/go-jsonnet v0.16.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ugorji/go v1.1.8 // indirect
	gitlab.com/king011/king-go v0.0.10
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201031054903-ff519b6c9102
	gopkg.in/yaml.v2 v2.3.0
	v2ray.com/core v0.0.0-00010101000000-000000000000
)

replace v2ray.com/core => ../v2ray-core
