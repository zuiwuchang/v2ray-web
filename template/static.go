package template

import (
	_ "embed"
)

//go:embed reader.js
var Default string

//go:embed reader_proxy.js
var Proxy string

//go:embed iptables_init.sh
var IPTablesInit string

//go:embed iptables_clear.sh
var IPTablesClear string
