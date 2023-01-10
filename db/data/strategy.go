package data

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(Strategy{})
}

const StrategyBucket = "strategy"
const StrategyDefault = "Default"

// 代理策略
type Strategy struct {
	// 唯一的名稱 供人類查看
	// 名稱 default 是系統保留的策略，其它策略將繼承這個策略 Name 和 Value 之外的所有值
	Name string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"`

	// 供腳本參考的 策略值 ，腳本應該使用此值生成 v2ray 的配置
	//
	//
	// 系統定義了幾個默認值，但如何處理它們完全是腳本決定的
	// * 0 使用默認的代理規則
	// * 1 全域代理
	// * 100 略區域網路的代理
	// * 200 略過區域網路和西朝鮮的代理
	// * 900 直連優先 (僅對非西朝鮮網路使用代理)
	// * 1000 直接連接
	Value int

	// 靜態 ip 列表
	// baidu.com 127.0.0.1
	// dns.google 8.8.8.8 8.8.4.4
	Host []Host

	// 這些 ip 使用代理
	ProxyIP []string
	// 這些 域名 使用代理
	ProxyDomain []string

	// 這些 ip 直連
	DirectIP []string
	// 這些 域名 直連
	DirectDomain []string

	// 這些 ip 阻止訪問
	BlockIP []string
	// 這些 域名 阻止訪問
	BlockDomain []string
}
type Host struct {
	Host string
	IP   []string
}

func (s *Strategy) Decode(b []byte) (e error) {
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	e = decoder.Decode(s)
	return
}

func (s *Strategy) Encoder() (b []byte, e error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	e = encoder.Encode(s)
	if e == nil {
		b = buffer.Bytes()
	}
	return
}
