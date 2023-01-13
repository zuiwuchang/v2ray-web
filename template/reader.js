// es5 以註釋 es5 開頭表示使用 es5 而非 golang template 創建配置，推薦使用 es5 比 golang 模板好寫很多

// 定義一個上下文，用於指導腳本如何工作，通常默認腳本就可以很好的工作了，你只需要修改這個上下文來自定義一些選項
const Context = {
    // 日誌 
    log: {
        // "debug" | "info" | "warning" | "error" | "none"
        level: "warning"
    },
    // socks5 代理設置
    socks5: {
        listen: "127.0.0.1",  // 監聽地址
        port: 11080,          // 監聽地址，不定義則不啓用
        udp: true,            // 是否代理 udp
        /** socks5 驗證用戶，爲空，則不驗證用戶*
        accounts: [
            {
                user: "fuck ccp",
                pass: "ccp go die",
            },
        ],/** */
    },
    // http 代理
    http: {
        listen: "127.0.0.1",  // 監聽地址
        port: 18118,          // 監聽地址，不定義則不啓用
        /** Basic Authentication，爲空，則不驗證用戶*
        accounts: [
            {
                user: "fuck ccp",
                pass: "ccp go die",
            },
        ],/** */
    },
    // tproxy 透明代理
    proxy: {
        // listen: "127.0.0.1",  // 監聽地址
        port: 22345,        // 監聽地址，不定義則不啓用
    },
}

// 這個函數將被系統調用
// 它返回一個 Object 將被轉爲 v2ray 的 JSON 配置
function render(opts) {
    const reader = new Reader(opts, Context)
    return reader.reader()
}
function getAnyObject(val) {
    if (typeof val === "object") {
        return val
    }
    return {}
}
function getAnyArray(val) {
    if (Array.isArray(val)) {
        return val
    }
    return []
}

function getAnyString(val) {
    if (typeof val === "string") {
        if (val == '') {
            return undefined
        }
        return val
    }
    return undefined
}
function getAnyInt(val) {
    if (typeof val === "number") {
        return val
    }
    return undefined
}
function getAnyArray(val) {
    if (Array.isArray(val)) {
        return val
    }
    return []
}
function getAnyArray2(val) {
    if (Array.isArray(val)) {
        const result = []
        for (let i = 0; i < val.length; i++) {
            result.push(getAnyArray(val[i]))
        }
        return result
    }
    return []
}
class Outbound {
    constructor(opts) {
        this.name = getAnyString(opts.Name)
        this.address = getAnyString(opts.Add)
        this.port = getAnyString(opts.Port)
        this.host = getAnyString(opts.Host)
        this.tls = getAnyString(opts.TLS)
        this.network = getAnyString(opts.Net)
        this.path = getAnyString(opts.Path)
        this.userID = getAnyString(opts.UserID)
        this.alterID = getAnyString(opts.AlterID)
        this.security = getAnyString(opts.Security)
        this.level = getAnyString(opts.Level)
        this.protocol = getAnyString(opts.Protocol)
        this.flow = getAnyString(opts.Flow)
    }
    toString() {
        return JSON.stringify(this, undefined, "\t")
    }
}
class Strategy {
    constructor(opts) {
        this.name = getAnyString(opts.Name)
        this.value = getAnyInt(opts.Value) ?? 900
        if (!Number.isSafeInteger(this.value) || this.value == 0) {
            this.value = 900
        }
        this.host = getAnyArray2(opts.Host)        // 靜態 ip 列表 [['baidu.com', '127.0.0.1'], ['dns.google', '8.8.8.8', '8.8.4.4']]
        this.proxyIP = getAnyArray(opts.ProxyIP)      // 這些 ip 使用代理
        this.proxyDomain = getAnyArray(opts.ProxyDomain)  // 這些 域名 使用代理
        this.directIP = getAnyArray(opts.DirectIP)     // 這些 ip 直連
        this.directDomain = getAnyArray(opts.DirectDomain) // 這些 域名 直連
        this.blockIP = getAnyArray(opts.BlockIP)      // 這些 ip 阻止訪問
        this.blockDomain = getAnyArray(opts.BlockDomain)  // 這些 域名 阻止訪問
    }
}
class ReaderOptions {
    constructor(opts) {
        this.path = getAnyString(opts.BasePath)
        this.ip = getAnyString(opts.AddIP)
        this.outbound = new Outbound(getAnyObject(opts.Outbound))
        this.strategy = new Strategy(getAnyObject(opts.Strategy))
    }
    toString() {
        return JSON.stringify(this, undefined, "\t")
    }
}
class ReaderContext {
    constructor(opts) {
        let obj = getAnyObject(opts.socks5)
        this.socks5 = {
            listen: getAnyString(obj.listen),
            port: getAnyInt(obj.port),
            udp: obj.udp ? true : false,
            accounts: getAnyArray(obj.accounts),
        }
        obj = getAnyObject(opts.http)
        this.http = {
            listen: getAnyString(obj.listen),
            port: getAnyInt(obj.port),
            accounts: getAnyArray(obj.accounts),
        }
        obj = getAnyObject(opts.proxy)
        this.proxy = {
            listen: getAnyString(obj.listen),
            port: getAnyInt(obj.port),
        }
        obj = getAnyObject(opts.log)
        this.log = {
            level: getAnyString(obj.level),
        }
    }
    toString() {
        return JSON.stringify(this, undefined, "\t")
    }
}
class Rule {
    domain = []
    ip = []
    _domain = new Set()
    _ip = new Set()

    pushDomain(a) {
        this._push(a)
        return this
    }
    pushIP(a) {
        this._push(a, true)
        return this
    }
    _push(a, ip) {
        if (!Array.isArray(a) || a.length == 0) {
            return
        }
        const keys = ip ? this._ip : this._domain
        const vals = ip ? this.ip : this.domain
        for (let i = 0; i < a.length; i++) {
            if (typeof a[i] !== "string") {
                continue
            }
            const val = a[i].trim()
            if (keys.has(val)) {
                continue
            }
            keys.add(val)
            vals.push(val)
        }
    }
}
class Reader {
    constructor(opts, ctx) {
        this.opts = new ReaderOptions(getAnyObject(opts))
        this.ctx = new ReaderContext(getAnyObject(ctx))
    }
    reader() {
        return {
            log: this.log(),
            dns: this.dns(),
            //     inbounds: renderInbounds(ctx),
            //     outbounds: renderOutbounds(ctx),
            //     routing: renderRouting(ctx),
        }
    }
    log() {
        const log = this.ctx.log
        return {
            // "debug" | "info" | "warning" | "error" | "none"
            loglevel: log.level,
        }
    }
    dns() {
        const opts = this.opts
        // 將代理服務器域名加入 靜態 dns
        const hosts = {}
        hosts[opts.outbound.address] = opts.ip
        const strategy = opts.strategy
        for (let i = 0; i < strategy.host.length; i++) {
            const host = strategy.host[i]
            if (host.length == 2) {
                hosts[host[0]] = host[1]
            } else {
                hosts[host[0]] = host.slice(1)
            }
        }

        if (strategy.value < 2) {// 全部代理
            return {
                hosts: hosts,
                servers: [
                    // 解析 非西朝 域名
                    "8.8.8.8", // google
                    "1.1.1.1", // cloudflare
                    "https+local://doh.dns.sb/dns-query",
                ],
            }
        }
        else if (strategy.value >= 1000) { // 全部直接連接
            return {
                hosts: hosts,
                servers: [
                    // 解析 西朝 域名
                    "119.29.29.29", // 騰訊 dns
                    "223.5.5.5", // 阿里 dns
                    "localhost",
                ],
            }
        }
        // 解析規則
        const proxy = new Rule()
            .pushDomain([
                "geosite:apple",
                "geosite:google",
                "geosite:microsoft",
                "geosite:facebook",
                "geosite:twitter",
                "geosite:telegram",
                "geosite:geolocation-!cn",
                "tld-!cn",
            ])
            .pushDomain(strategy.proxyDomain)
            .pushIP(strategy.proxyIP)
        const direct = new Rule()
            .pushDomain([
                "geosite:cn",
            ])
            .pushIP([
                "geoip:cn",
            ])
            .pushDomain(strategy.directDomain)
            .pushIP(strategy.directIP)

        if (strategy.value < 900) { // 代理優先
            return {
                hosts: hosts,
                servers: [
                    // 解析 非西朝 域名
                    {
                        address: "8.8.8.8", // google
                        port: 53,
                        domains: proxy.domain,
                        expectIPs: proxy.ip,
                    },
                    {
                        address: "1.1.1.1", // cloudflare
                        port: 53,
                        domains: proxy.domain,
                        expectIPs: proxy.ip,
                    },
                    // 解析 西朝 域名
                    {
                        address: "119.29.29.29", // 騰訊 dns
                        port: 53,
                        domains: direct.domain,
                        expectIPs: direct.ip,
                    },
                    {
                        address: "223.5.5.5", // 阿里 dns
                        port: 53,
                        domains: direct.domain,
                        expectIPs: direct.ip,
                    },
                    // 未匹配的
                    "8.8.8.8",
                    "1.1.1.1",
                    "https+local://doh.dns.sb/dns-query"
                ],
            }
        }
        // 直連優先
        return {
            hosts: hosts,
            servers: [
                // 解析 西朝 域名
                {
                    address: "119.29.29.29", // 騰訊 dns
                    port: 53,
                    domains: direct.domain,
                    expectIPs: direct.ip,
                },
                {
                    address: "223.5.5.5", // 阿里 dns
                    port: 53,
                    domains: direct.domain,
                    expectIPs: direct.ip,
                },
                // 解析 非西朝 域名
                {
                    address: "8.8.8.8", // google
                    port: 53,
                    domains: proxy.domain,
                    expectIPs: proxy.ip,
                },
                {
                    address: "1.1.1.1", // cloudflare
                    port: 53,
                    domains: proxy.domain,
                    expectIPs: proxy.ip,
                },
                // 未匹配的
                "119.29.29.29", // 騰訊 dns
                "223.5.5.5", // 阿里 dns
                "localhost",
            ],
        }



    }
}
function getStrategy(ctx) {
    const strategy = ctx.Strategy
    if (typeof strategy !== "object") {
        return {
            Value: 900,
        }
    }
    let value = 900
    if (typeof strategy.Value === "number" && strategy.value != 0) {
        value = strategy.value
    }

    return {
        Value: value,
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
                address: "119.29.29.29", // 騰訊 dns
                port: 53,
                domains: ["geosite:cn"],
                expectIPs: ["geoip:cn"]
            },
            {
                address: "223.5.5.5", // 阿里 dns
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
function getListen(addr) {
    if (typeof addr === "string") {
        addr = addr.trim()
        if (addr != "") {
            return addr
        }
    }
    return undefined
}
function isValidPort(p) {
    return Number.isSafeInteger(p) && p > 0 && p <= 65535
}
function renderInbounds(ctx) {
    const result = []
    const socks5 = Context.socks5
    // 本地 socks5 代理
    if (socks5 && isValidPort(socks5.port)) {
        const accounts = socks5.accounts
        const password = Array.isArray(accounts) && accounts.length > 0
        result.push({
            tag: "socks",
            protocol: "socks",
            listen: getListen(socks5.listen),
            port: socks5.port,
            settings: {
                auth: password ? "password" : "noauth",
                accounts: password ? accounts : undefined,
                udp: socks5.udp ? true : false,
                userLevel: 0,
            }
        })
    }
    // 本地 http 代理
    const http = Context.http
    if (http && isValidPort(http.port)) {
        const accounts = http.accounts
        const password = Array.isArray(accounts) && accounts.length > 0
        result.push({
            tag: "http",
            protocol: "http",
            listen: getListen(http.listen),
            port: http.port,
            timeout: 300, // 超時時間爲 300 秒
            allowTransparent: false, // 爲 true 不止代理也轉發所有 http 請求
            accounts: password ? accounts : undefined,
            userLevel: 0,
        })
    }
    // 透明代理
    const proxy = Context.proxy
    if (proxy && isValidPort(proxy.port)) {
        result.push({
            tag: "all-in",
            protocol: "dokodemo-door",
            listen: getListen(proxy.listen),
            port: proxy.port,
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
        })
    }
    return result

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