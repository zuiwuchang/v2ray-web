// es5 以註釋 es5 開頭表示使用 es5 而非 golang template 創建配置，推薦使用 es5 比 golang 模板好寫很多

// 這個函數將被系統調用
// 它返回一個 Object 將被轉爲 v2ray 的 JSON 配置
function render(ctx) {
    return {
        log: renderLog(ctx),
        dns: renderDNS(ctx),
        inbounds: renderInbounds(ctx),
        outbounds: renderOutbounds(ctx),
        routing: renderRouting(ctx),
    }
}

function renderLog(ctx) {
    return {
        // "debug" | "info" | "warning" | "error" | "none"
        loglevel: "warning",
    }
}
function renderDNS(ctx) {
    // 將代理服務器域名加入 靜態 dns
    const hosts = {}
    hosts[ctx.Outbound.Add] = ctx.AddIP
    return {
        hosts: hosts,
        servers: [
            // 解析 西朝 域名
            {
                address: "119.29.29.29",
                port: 53,
                domains: ["geosite:cn"],
                expectIPs: ["geoip:cn"]
            },
            {
                address: "223.5.5.5",
                port: 53,
                domains: ["geosite:cn"],
                expectIPs: ["geoip:cn"]
            },
            // 解析 非西朝 域名
            "8.8.8.8",
            "1.1.1.1",
            "https+local://doh.dns.sb/dns-query"
        ],
    }
}
function renderInbounds(ctx) {
    return [
        // 本地 socks5 代理
        {
            tag: "socks",
            listen: "127.0.0.1",
            protocol: "socks",
            port: 1080,
            settings: {
                auth: "noauth",
                accounts: [ // auth 爲 "password" 有效，配置 socks5 用戶名和密碼
                    {
                        user: "fuck ccp",
                        pass: "ccp go die",
                    },
                ],
                udp: true,
                userLevel: 0,
            }
        },
        // 本地 http 代理
        {
            tag: "http",
            listen: "127.0.0.1",
            protocol: "http",
            port: 8118,
            timeout: 300, // 超時時間爲 300 秒
            allowTransparent: false, // 爲 true 不止代理也轉發所有 http 請求
            accounts: [ // 非空則要求代理設置 Basic Authentication
                /**
                {
                    user: "fuck ccp",
                    pass: "ccp go die",
                },
                /**/
            ],
            userLevel: 0,
        },
        // 透明代理
        {
            tag: "all-in",
            protocol: "dokodemo-door",
            port: 12345,
            settings: {
                network: "tcp,udp",
                followRedirect: true
            },
            sniffing: {
                enabled: true,
                destOverride: [
                    "http",
                    "tls",
                ],
            },
            streamSettings: {
                sockopt: {
                    tproxy: "tproxy"
                }
            }
        },
    ]
}
function intValue(val, def) {
    if (val == "") {
        return def
    } else if (typeof val === "number") {
        return val
    } else if (typeof val === "string") {
        const v = parseInt(val)
        if (isFinite(val)) {
            return v
        }
    }
    return def
}

function xtlsFlow(ctx) {
    if (ctx.Outbound.TLS != "xtls") {
        return
    } else if (ctx.Outbound.Protocol != "vless" && ctx.Outbound.Protocol != "trojan") {
        return
    }
    const flow = ctx.Outbound.Flow
    if (typeof flow === "string" && flow != "") {
        return flow
    }
}
function tlsSettings(ctx) {
    return {
        serverName: ctx.Outbound.Host == '' ? ctx.Outbound.Add : ctx.Outbound.Host,
        rejectUnknownSni: false,
        alpn: ["h2", "http/1.1"],
        // allowInsecure: true,//允許不安全的證書
        // "" | "chrome" | "firefox" | "safari" | "randomized"
        fingerprint: "firefox", // 模擬 tls 指紋
    }
}
function tcpSettings(ctx) {
    return {
        header: {
            type: "none",
        },
    }
}
function kcpSettings(ctx) {
    return {
        mtu: 1350,
        tti: 20,
        uplinkCapacity: 5,
        downlinkCapacity: 20,
        congestion: false,
        readBufferSize: 1,
        writeBufferSize: 1,
        header: {
            type: "none"
        },
    }
}
function wsSettings(ctx) {
    return {
        path: ctx.Outbound.Path == "" ? "/" : ctx.Outbound.Path,
        headers: {
            Host: ctx.Outbound.Host == '' ? ctx.Outbound.Add : ctx.Outbound.Host,
        },
    }
}
function httpSettings(ctx) {
    return {
        path: ctx.Outbound.Path == "" ? "/" : ctx.Outbound.Path,
        method: "PUT",
        headers: {
            Host: ctx.Outbound.Host == '' ? ctx.Outbound.Add : ctx.Outbound.Host,
        },
    }
}
function dsSettings(ctx) {
    return {
        path: ctx.Outbound.Path == "" ? "/" : ctx.Outbound.Path,
        abstract: false,
        padding: false,
    }
}
function quicSettings(ctx) {
    return {
        security: "none",
        key: "",
        header: {
            type: "none"
        }
    }
}
function sockopt(ctx) {
    return {
        mark: 2,
    }
}
function renderOutbounds(ctx) {
    let proxy
    switch (ctx.Outbound.Protocol) {
        case "shadowsocks":
            proxy = outboundsShadowsocks(ctx)
            break;
        case "trojan":
            proxy = outboundsTrojan(ctx)
            break
        case "vless":
            proxy = outboundsVless(ctx)
            break
        case "vmess":
            proxy = outboundsVmess(ctx)
            break
        default:
            throw "not support protocol: " + ctx.Outbound.Protocol
    }
    proxy.tag = "proxy"
    proxy.protocol = ctx.Outbound.Protocol
    return [
        // 直接 訪問
        {
            tag: "direct",
            protocol: "freedom",
            streamSettings: {
                sockopt: {
                    mark: 2,
                },
            },
        },
        // 代理服務器
        proxy,
        // 拒絕 訪問
        {
            tag: "block",
            protocol: "blackhole",
            settings: {}
        },
        {
            tag: "dns-out",
            protocol: "dns",
            settings: {
                address: "8.8.8.8"
            },
            proxySettings: {
                tag: "proxy"
            },
            streamSettings: {
                sockopt: {
                    mark: 2,
                },
            },
        },
    ]
}
function outboundsShadowsocks(ctx) {
    return {
        settings: {
            servers: [
                {
                    address: ctx.AddIP,
                    port: intValue(ctx.Outbound.Port),
                    password: ctx.Outbound.UserID,
                    method: ctx.Outbound.Security,
                    uot: false,
                    level: intValue(ctx.Outbound.Level, 0),
                },
            ],
        },
        streamSettings: {
            sockopt: sockopt(ctx),
        },
    }
}
function outboundsTrojan(ctx) {
    const xtls = ctx.Outbound.TLS == "xtls"
    return {
        settings: {
            servers: [
                {
                    address: ctx.AddIP,
                    port: intValue(ctx.Outbound.Port),
                    password: ctx.Outbound.UserID,
                    flow: xtlsFlow(ctx),
                    level: intValue(ctx.Outbound.Level, 0),
                },
            ],
        },
        streamSettings: {
            network: "tcp",
            security: xtls ? "xtls" : "tls",
            tlsSettings: xtls ? undefined : tlsSettings(ctx),
            xtlsSettings: xtls ? tlsSettings(ctx) : undefined,
            sockopt: sockopt(ctx),
        },
    }
}
function outboundsVless(ctx) {
    return {
        settings: {
            vnext: [
                {
                    address: ctx.AddIP,
                    port: intValue(ctx.Outbound.Port),
                    users: [
                        {
                            id: ctx.Outbound.UserID,
                            flow: xtlsFlow(ctx),
                            encryption: "none",
                            level: intValue(ctx.Outbound.Level, 0),
                        },
                    ],
                },
            ],
        },
        streamSettings: {
            network: ctx.Outbound.Net,
            security: ctx.Outbound.TLS,
            tlsSettings: ctx.Outbound.TLS == "tls" ? tlsSettings(ctx) : undefined,
            xtlsSettings: ctx.Outbound.TLS == "xtls" ? tlsSettings(ctx) : undefined,
            tcpSettings: ctx.Outbound.Net == "tcp" ? tcpSettings(ctx) : undefined,
            kcpSettings: ctx.Outbound.Net == "kcp" ? kcpSettings(ctx) : undefined,
            wsSettings: ctx.Outbound.Net == "ws" ? wsSettings(ctx) : undefined,
            httpSettings: ctx.Outbound.Net == "http" ? httpSettings(ctx) : undefined,
            dsSettings: ctx.Outbound.Net == "domainsocket" ? dsSettings(ctx) : undefined,
            quicSettings: ctx.Outbound.Net == "quic" ? quicSettings(ctx) : undefined,
            sockopt: sockopt(ctx),
        },
    }
}
function outboundsVmess(ctx) {
    return {
        settings: {
            vnext: [
                {
                    address: ctx.AddIP,
                    port: intValue(ctx.Outbound.Port),
                    users: [
                        {
                            id: ctx.Outbound.UserID,
                            alterId: intValue(ctx.Outbound.AlterID, 0),
                            security: ctx.Outbound.Security,
                            level: intValue(ctx.Outbound.Level, 0),
                        },
                    ],
                },
            ],
        },
        streamSettings: {
            network: ctx.Outbound.Net,
            security: ctx.Outbound.TLS,
            tlsSettings: ctx.Outbound.TLS == "tls" ? tlsSettings(ctx) : undefined,
            xtlsSettings: ctx.Outbound.TLS == "xtls" ? tlsSettings(ctx) : undefined,
            tcpSettings: ctx.Outbound.Net == "tcp" ? tcpSettings(ctx) : undefined,
            kcpSettings: ctx.Outbound.Net == "kcp" ? kcpSettings(ctx) : undefined,
            wsSettings: ctx.Outbound.Net == "ws" ? wsSettings(ctx) : undefined,
            httpSettings: ctx.Outbound.Net == "http" ? httpSettings(ctx) : undefined,
            dsSettings: ctx.Outbound.Net == "domainsocket" ? dsSettings(ctx) : undefined,
            quicSettings: ctx.Outbound.Net == "quic" ? quicSettings(ctx) : undefined,
            sockopt: sockopt(ctx),
        },
    }
}

function renderRouting(ctx) {
    return {
        domainStrategy: "IPIfNonMatch",
        rules: [
            // 攔截域名解析
            {
                type: "field",
                inboundTag: ["all-in"],
                port: 53,
                outboundTag: "dns-out"
            },
            {
                type: "field",
                ip: ["8.8.8.8", "1.1.1.1"],
                outboundTag: "proxy"
            },
            /* 屏蔽廣告 *
            {
                type: "field",
                domain: ["geosite:category-ads-all"],
                // domain: ["geosite:category-ads"],
                outboundTag: "block"
            },/* */
            // 代理 非西朝 域名
            {
                type: "field",
                domain: ["geosite:geolocation-!cn"],
                outboundTag: "proxy"
            },
            {
                type: "field",
                ip: ["geoip:telegram"],
                outboundTag: "proxy"
            },
            // 不代理 bt 下載
            {
                type: "field",
                protocol: ["bittorrent"],
                outboundTag: "direct"
            },
            // 不代理 西朝
            {
                type: "field",
                ip: ["geoip:private", "geoip:cn"],
                outboundTag: "direct"
            },
            {
                type: "field",
                domain: ["geosite:cn"],
                outboundTag: "direct"
            },
        ],
    }
}