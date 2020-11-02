package speed

const templateText = `{
    "log": {
        "loglevel": "none"
    },
    "inbounds": [
        {
            "tag": "socks",
            "listen": "127.0.0.1",
            "protocol": "socks",
            "port": %v,
            "settings": {
                "auth": "noauth"
            }
        }
    ],
    "outbounds": [
        {
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
        }
    ]
}
`
