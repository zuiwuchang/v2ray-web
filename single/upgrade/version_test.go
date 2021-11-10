package upgrade_test

import (
	"testing"

	"gitlab.com/king011/v2ray-web/single/upgrade"
)

func TestVersion(t *testing.T) {
	v := upgrade.ParseVersion(`v1.3.1-9-g022095a`)
	if v.X != 1 || v.Y != 3 || v.Z != 1 {
		t.Fatal(`ParseVersion err`, v)
	}
	v0 := upgrade.ParseVersion(`v1.3.2`)
	if !v.LessMatch(&v0) {
		t.Fatal(`not match`)
	}
	v0 = upgrade.ParseVersion(`v2.3.2`)
	if v.LessMatch(&v0) {
		t.Fatal(`not match`)
	}

	v0 = upgrade.ParseVersion(`v0.3.2`)
	if !v0.LessMatch(&v) {
		t.Fatal(`not match`)
	}
}
