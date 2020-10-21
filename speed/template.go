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
            "protocol": "{{.Protocol}}",
            "settings": {
                "vnext": [
                    {{.Vnext}}
                ]
            },
            "streamSettings": {{.StreamSettings}}
        }
    ]
}
`
