export const ContextText = `{
    BasePath string
    AddIP string
    Outbound {
      Name string 
      Add string
      Port string
      Host string
      TLS string
      Net string
      Path string
      UserID string
      AlterID string
      Security string
      Level string
      Protocol string
    }
  }`
export const V2rayTemplate = `{
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
        // 本地 http 代理
        {
          "tag": "http",
          "listen": "127.0.0.1",
          "protocol": "http",
          "port": 8118
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
            "protocol": "{{.Outbound.Protocol}}",
  {{if eq .Outbound.Protocol "shadowsocks"}}
  "settings": {
      "servers": [
          {
            "address": "{{.Outbound.Add}}",
            "port": {{.Outbound.Port}},
            "method": "{{.Outbound.Security}}",
            "password": "{{.Outbound.UserID}}",
            "ota": false,
            "level": 0
          }
      ]
  }
  {{else}}
  "settings": {
      "vnext": [{
          "address": "{{.Outbound.Add}}",
          "port": {{.Outbound.Port}},
          "users": [{
              "id": "{{.Outbound.UserID}}",
  
               {{if eq .Outbound.Protocol "vmess"}}
                  {{if eq .Outbound.AlterID ""}}
                          "alterId": 0,
                  {{else}}
                          "alterId": {{.Outbound.AlterID}},
                  {{end}}
                  "security": "{{.Outbound.Security}}",
              {{else if eq .Outbound.Protocol "vless"}}
                  "flow": "",
                  "encryption": "none",
              {{end}}
  
              {{if eq .Outbound.Level ""}}
                      "level": 0
              {{else}}
                      "level": {{.Outbound.Level}}
              {{end}}
          }]
      }]
  },
  "streamSettings": {
      "network": "{{.Outbound.Net}}",
      "security": "{{.Outbound.TLS}}",
      {{if eq .Outbound.TLS "tls"}}
          "tlsSettings": {
              "serverName": "{{.Outbound.Host}}",
              "allowInsecure": false,
              "alpn": ["http/1.1"],
              "certificates": [],
              "disableSystemRoot": false
          },
      {{end}}
  
      {{if eq .Outbound.Net "tcp"}}
          "tcpSettings": {
              "header": {
                  "type": "none"
              }
          },
      {{else if eq .Outbound.Net "kcp"}}
          "kcpSettings": {
              "mtu": 1350,
              "tti": 20,
              "uplinkCapacity": 5,
              "downlinkCapacity": 20,
              "congestion": false,
              "readBufferSize": 1,
              "writeBufferSize": 1,
              "header": {
                  "type": "none"
              }
          },
      {{else if eq .Outbound.Net "ws"}}
          "wsSettings": {
              {{if eq .Outbound.Path ""}}
                  "path": "/",
              {{else}}
                  "path": "{{.Outbound.Path}}",
              {{end}}
              
              
              "headers": {
                  {{if eq .Outbound.Host ""}}
                  {{else}}
                      "Host": "{{.Outbound.Host}}"
                  {{end}}
              }
          },
      {{else if eq .Outbound.Net "http"}}
          "httpSettings": {
              {{if eq .Outbound.Host ""}}
              {{else}}
                  "Host": ["{{.Outbound.Host}}"]
              {{end}}
  
              {{if eq .Outbound.Path ""}}
                  "path": "/"
              {{else}}
                  "path": "{{.Outbound.Path}}"
              {{end}}
          },
      {{else if eq .Outbound.Net "domainsocket"}}
          "dsSettings": {
              {{if eq .Outbound.Path ""}}
                  "path": "/"
              {{else}}
                  "path": "{{.Outbound.Path}}"
              {{end}}
          },
      {{else if eq .Outbound.Net "quic"}}
          "quicSettings": {
              "security": "none",
              "key": "",
              "header": {
                  "type": "none"
              }
          },
      {{end}}
      "sockopt": {
          "mark": 0,
          "tcpFastOpen": false,
          "tproxy": "off"
      }
  },
  "mux": {
      "enabled": true
  }
  {{end}}
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

