package net_test

import (
	"testing"

	"gitlab.com/king011/v2ray-web/internal/net"
)

func testValue(t *testing.T, name, protocol string, expect, current *net.Outbound) {
	if name != protocol {
		t.Fatalf(`protocol: expect "%s", but get "%s".`, name, protocol)
	}

	if expect.Name != current.Name {
		t.Fatalf(`%s.Name: expect "%s", but get "%s".`, name, expect.Name, current.Name)
	}
	if expect.Add != current.Add {
		t.Fatalf(`%s.Add: expect "%s", but get "%s".`, name, expect.Add, current.Add)
	}
	if expect.Port != current.Port {
		t.Fatalf(`%s.Port: expect "%s", but get "%s".`, name, expect.Port, current.Port)
	}
	if expect.Host != current.Host {
		t.Fatalf(`%s.Host: expect "%s", but get "%s".`, name, expect.Host, current.Host)
	}
	if expect.TLS != current.TLS {
		t.Fatalf(`%s.TLS: expect "%s", but get "%s".`, name, expect.TLS, current.TLS)
	}
	if expect.Net != current.Net {
		t.Fatalf(`%s.Net: expect "%s", but get "%s".`, name, expect.Net, current.Net)
	}
	if expect.Path != current.Path {
		t.Fatalf(`%s.Path: expect "%s", but get "%s".`, name, expect.Path, current.Path)
	}
	if expect.UserID != current.UserID {
		t.Fatalf(`%s.UserID: expect "%s", but get "%s".`, name, expect.UserID, current.UserID)
	}
	if expect.AlterID != current.AlterID {
		t.Fatalf(`%s.AlterID: expect "%s", but get "%s".`, name, expect.AlterID, current.AlterID)
	}
	if expect.Security != current.Security {
		t.Fatalf(`%s.Security: expect "%s", but get "%s".`, name, expect.Security, current.Security)
	}
	if expect.Level != current.Level {
		t.Fatalf(`%s.Level: expect "%s", but get "%s".`, name, expect.Level, current.Level)
	}
}
func TestTrojan(t *testing.T) {
	expect := &net.Outbound{
		Name:   `測試`,
		UserID: `userid`,
		Add:    `hostname`,
		Port:   `22584`,
		Level:  `2`,
	}
	protocol, outbound := net.AnalyzeString("trojan://userid@hostname:22584?name=%E6%B8%AC%E8%A9%A6&level=2")
	testValue(t, `trojan`, protocol, expect, outbound)
	protocol, outbound = net.AnalyzeString("trojan://userid@hostname:22584?name=測試&level=2")
	testValue(t, `trojan`, protocol, expect, outbound)

	protocol, outbound = net.AnalyzeString("trojan://userid@hostname:22584?level=2#%E6%B8%AC%E8%A9%A6")
	testValue(t, `trojan`, protocol, expect, outbound)
	protocol, outbound = net.AnalyzeString("trojan://userid@hostname:22584?level=2#測試")
	testValue(t, `trojan`, protocol, expect, outbound)
}
