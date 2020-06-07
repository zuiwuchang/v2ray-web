module gitlab.com/king011/v2ray-web

go 1.12

require (
	github.com/boltdb/bolt v1.3.1
	github.com/gin-contrib/gzip v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-jsonnet v0.16.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.9
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	gitlab.com/king011/king-go v0.0.10
	go.starlark.net v0.0.0-20200519165436-0aa95694c768 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	google.golang.org/genproto v0.0.0-20200605102947-12044bf5ea91 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	v2ray.com/core v0.0.0-00010101000000-000000000000
	v2ray.com/ext v4.15.0+incompatible
)

replace v2ray.com/core => github.com/v2ray/v2ray-core v4.23.4+incompatible
