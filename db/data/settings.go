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

// V2rayTemplate v2ray 默認設定模板
const V2rayTemplate = `{
    "log": {
        "loglevel": "warning"
    },
    "dns": {
        "servers": [
            // 使用 google 解析
            {
                "address": "8.8.8.8",
                "port": 53,
                "domains": [
                    "geosite:google",
                    "geosite:facebook",
                    "geosite:geolocation-!cn"
                ]
            },
            // 使用 趙國 解析服務
            {
                "address": "114.114.114.114",
                "port": 53,
                "domains": [
                    "geosite:cn",
                    "geosite:speedtest",
                    "domain:cn"
                ]
            },
            "8.8.8.8",
            "8.8.4.4",
            "localhost"
        ]
    },
    "inbounds": [
        // 本地 socks5 代理
        {
            "tag": "socks",
            "listen": "127.0.0.1",
            "protocol": "socks",
            "port": 1080,
            "settings": {
                "auth": "noauth"
            }
        },
        // 透明代理
        {
            "tag": "redir",
            "protocol": "dokodemo-door",
            "port": 10090,
            "settings": {
                "network": "tcp,udp",
                "followRedirect": true
            },
            "sniffing": {
                "enabled": true,
                "destOverride": [
                    "http",
                    "tls"
                ]
            }
        },
        // dns 代理 解決 域名污染
        {
            "tag": "dns",
            "protocol": "dokodemo-door",
            "port": 10054,
            "settings": {
                "address": "8.8.8.8",
                "port": 53,
                "network": "tcp,udp",
                "followRedirect": false
            }
        }
    ],
    "outbounds": [
        // 代理 訪問
        {
            "tag": "proxy",
            "protocol": "vmess",
            "settings": {
                "vnext": [
                    {{.Vnext}}
                ]
            },
            "streamSettings": {{.StreamSettings}},
            "mux": {
                "enabled": true
            }
        },
        // 直接 訪問
        {
            "tag": "freedom",
            "protocol": "freedom",
            "settings": {}
        },
        // 拒絕 訪問
        {
            "tag": "blackhole",
            "protocol": "blackhole",
            "settings": {}
        }
    ],
    "routing": {
        "domainStrategy": "IPIfNonMatch",
        "rules": [
            // 通過透明代理 進入 一律 代理訪問
            {
                "type": "field",
                "network": "tcp,udp",
                "inboundTag": [
                    "redir",
                    "dns"
                ],
                "outboundTag": "proxy"
            },
            // 代理訪問
            {
                "type": "field",
                "domain": [
                    "geosite:google",
                    "geosite:facebook",
                    "geosite:geolocation-!cn"
                ],
                "network": "tcp,udp",
                "outboundTag": "proxy"
            },
            // 直接訪問
            {
                "type": "field",
                "domain": [
                    "geosite:cn",
                    "geosite:speedtest",
                    "domain:cn",
                    "geoip:private"
                ],
                "network": "tcp,udp",
                "outboundTag": "freedom"
            },
            {
                "type": "field",
                "ip": [
                    "geoip:cn"
                ],
                "network": "tcp,udp",
                "outboundTag": "freedom"
            }
        ]
    }
}`

func init() {
	gob.Register(Settings{})
	gob.Register(IPTables{})
}

// Settings 系統設定
type Settings struct {
	URL      string `json:"url,omitempty"`
	V2ray    bool   `json:"v2ray,omitempty"`
	IPTables bool   `json:"iptables,omitempty"`
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
	Shell string `json:"shell,omitempty"`
	View  string `json:"view,omitempty"`
	Clear string `json:"clear,omitempty"`
	Init  string `json:"init,omitempty"`
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
	iptables.Clear = `iptables -t nat -F
iptables -t filter -F
iptables -t mangle -F`
	iptables.Init = `# 本地 dns 端口
DNS_Port=10053

# Redir 程式 端口
Redir_Port=10090

# 要放行的 ip 數組 
# 通常是 服務器地址 和 不需要代理的地址
IP_Servers=(
    {{.AddIP}}
    114.114.114.114
)

# 定義 內網 地址
# 一般不用修改
IP_Private=(
	0/8
	127/8
	10/8
	169.254/16
	172.16/12
	192.168/16
	224/4
	240/4
)

# 創建 nat/tcp 轉發鏈 用於 轉發 tcp流
iptables-save | egrep "^\:NAT_TCP" >> /dev/null
if [[ $? != 0 ]];then
    iptables -t nat -N NAT_TCP
fi

# 放行所有 內網地址
for i in ${!IP_Private[@]}
do
    iptables -t nat -A NAT_TCP -d ${IP_Private[i]} -j RETURN
done

# 放行 發往 服務器的 數據
for i in ${!IP_Servers[@]}
do
    iptables -t nat -A NAT_TCP -d ${IP_Servers[i]} -j RETURN
done

# 重定向 tcp 數據包到 redir 監聽端口
iptables -t nat -A NAT_TCP -p tcp -j REDIRECT --to-ports $Redir_Port

# 重定向 向網關發送的 dns 查詢
for i in ${!IP_Private[@]}
do
        iptables -t nat -A OUTPUT -d ${IP_Private[i]} -p udp -m udp --dport 53 -j DNAT --to-destination 127.0.0.1:$DNS_Port
        iptables -t nat -A OUTPUT -d ${IP_Private[i]} -p tcp -m tcp --dport 53 -j DNAT --to-destination 127.0.0.1:$DNS_Port
done

# 重定向 數據流向 NAT_TCP
iptables -t nat -A OUTPUT -p tcp -j NAT_TCP
iptables -t nat -A PREROUTING -p tcp -s 192.168/16 -j NAT_TCP
iptables -t nat -A POSTROUTING -s 192.168/16 -j MASQUERADE
`
}
