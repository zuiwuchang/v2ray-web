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
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/refraction-networking/utls v0.0.0-20200820030103-33a29038e742 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ugorji/go v1.1.8 // indirect
	gitlab.com/king011/king-go v0.0.10
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20200923182212-328152dc79b1
	golang.org/x/sys v0.0.0-20200923182605-d9f96fdee20d // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200923140941-5646d36feee1 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	v2ray.com/core v0.0.0-00010101000000-000000000000
)

replace v2ray.com/core => ../v2ray-core
