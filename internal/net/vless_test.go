package net_test

import (
	"testing"

	"gitlab.com/king011/v2ray-web/internal/net"
)

func TestVless(t *testing.T) {
	expect := &net.Outbound{
		Name:   `vless h2`,
		UserID: `userid`,
		Add:    `hostname`,
		Port:   `22583`,
		Level:  `1`,
		TLS:    `tls`,
		Net:    `http`,
		Path:   `/h2`,
	}
	protocol, outbound := net.AnalyzeString("vless://userid@hostname:22583?host=&security=tls&type=http&path=/h2&level=1#vless%20h2")
	testValue(t, `vless`, protocol, expect, outbound)
}
