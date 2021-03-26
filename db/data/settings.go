package data

import (
	"bytes"
	"encoding/gob"
)

const (
	// SettingsBucket .
	SettingsBucket = "settings"

	// SettingsV2ray v2ray 配置模板
	SettingsV2ray = "v2ray"

	// SettingsSettings 系統設定
	SettingsSettings = "settings"

	// SettingsIPTables iptables 防火牆 命令模板
	SettingsIPTables = "iptables"

	// SettingsLast 最後啓動的 v2ray 服務
	SettingsLast = "last"
)

func init() {
	gob.Register(Settings{})
	gob.Register(IPTables{})
}

// Settings 系統設定
type Settings struct {
	URL      string `json:"url,omitempty" json:"xml,omitempty" json:"yaml,omitempty"`
	V2ray    bool   `json:"v2ray,omitempty" xml:"v2ray,omitempty" yaml:"v2ray,omitempty"`
	IPTables bool   `json:"iptables,omitempty" xml:"iptables,omitempty" yaml:"iptables,omitempty"`
}

// Decode 由 []byte 解碼
func (settings *Settings) Decode(b []byte) (e error) {
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	e = decoder.Decode(settings)
	return
}

// Encoder 編碼到 []byte
func (settings *Settings) Encoder() (b []byte, e error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	e = encoder.Encode(settings)
	if e == nil {
		b = buffer.Bytes()
	}
	return
}

// ResetDefault 重新 置爲默認值
func (settings *Settings) ResetDefault() {
	settings.V2ray = true
	settings.IPTables = false
	settings.URL = `https://www.youtube.com/`
}

// IPTables 防火牆設置
type IPTables struct {
	Shell string `json:"shell,omitempty" xml:"shell,omitempty" yaml:"shell,omitempty"`
	View  string `json:"view,omitempty" xml:"view,omitempty" yaml:"view,omitempty"`
	Clear string `json:"clear,omitempty" xml:"clear,omitempty" yaml:"clear,omitempty"`
	Init  string `json:"init,omitempty" xml:"init,omitempty" yaml:"init,omitempty"`
}

// ResetDefault 重新 置爲默認值
func (iptables *IPTables) ResetDefault() {
	iptables.resetDefaultLinux()
}

// Decode 由 []byte 解碼
func (iptables *IPTables) Decode(b []byte) (e error) {
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	e = decoder.Decode(iptables)
	return
}

// Encoder 編碼到 []byte
func (iptables *IPTables) Encoder() (b []byte, e error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	e = encoder.Encode(iptables)
	if e == nil {
		b = buffer.Bytes()
	}
	return
}
func (iptables *IPTables) resetDefaultLinux() {
	iptables.Shell = "bash"
	iptables.View = "iptables-save"
	iptables.Clear = iptablesClear
	iptables.Init = iptablesInit
}
