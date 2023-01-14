// es5 以註釋 es5 開頭表示使用 es5 而非 golang template 創建配置，推薦使用 es5 比 golang 模板好寫很多

// 定義一個上下文，用於指導腳本如何工作，通常默認腳本就可以很好的工作了，你只需要修改這個上下文來自定義一些選項
const Context = {
    // 日誌 
    log: {
        // "debug" | "info" | "warning" | "error" | "none"
        level: "warning"
    },
    // 如何處理 bt, "direct" | "proxy" | undefined
    bittorrent: "direct",
    // socks5 代理設置
    socks5: {
        listen: "127.0.0.1",  // 監聽地址
        port: 1080,          // 監聽地址，不定義則不啓用
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
        port: 8118,          // 監聽地址，不定義則不啓用
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
        port: 12345,        // 監聽地址，不定義則不啓用
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
function toAnyInt(val) {
    if (typeof val === "string") {
        const v = parseInt(val)
        if (Number.isSafeInteger(v)) {
            return v
        }
        return undefined
    }
    if (typeof val === "number") {
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
        this.port = toAnyInt(opts.Port)
        this.host = getAnyString(opts.Host)
        this.tls = getAnyString(opts.TLS)
        this.network = getAnyString(opts.Net)
        this.path = getAnyString(opts.Path)
        this.userID = getAnyString(opts.UserID)
        this.alterID = toAnyInt(opts.AlterID)
        this.security = getAnyString(opts.Security)
        this.level = toAnyInt(opts.Level)
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
    toString() {
        return JSON.stringify(this, undefined, "\t")
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
        let port = getAnyInt(obj.port) ?? 0
        this.socks5 = port > 0 && port <= 65535 ? {
            listen: getAnyString(obj.listen),
            port: port,
            udp: obj.udp ? true : false,
            accounts: getAnyArray(obj.accounts),
        } : undefined

        obj = getAnyObject(opts.http)
        port = getAnyInt(obj.port) ?? 0
        this.http = port > 0 && port <= 65535 ? {
            listen: getAnyString(obj.listen),
            port: port,
            accounts: getAnyArray(obj.accounts),
        } : undefined

        obj = getAnyObject(opts.proxy)
        port = getAnyInt(obj.port) ?? 0
        this.proxy = port > 0 && port <= 65535 ? {
            listen: getAnyString(obj.listen),
            port: getAnyInt(obj.port),
        } : undefined

        obj = getAnyObject(opts.log)
        this.log = {
            level: getAnyString(obj.level),
        }

        this.bittorrent = getAnyString(opts.bittorrent)
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
    isValid() {
        return this.domain.length != 0 || this.ip.length != 0
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
            inbounds: this.inbounds(),
            outbounds: this.outbounds(),
            routing: this.routing(),
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

        if (strategy.value < 100) {// 全部代理
            // 解析規則
            const proxy = new Rule()
                .pushDomain(strategy.proxyDomain)
                .pushIP(strategy.proxyIP)
            const direct = new Rule()
                .pushDomain(strategy.directDomain)
                .pushIP(strategy.directIP)
            return this._dns(hosts, proxy, direct, false)
        } else if (strategy.value < 200) { // 代理公有 ip
            // 解析規則
            const proxy = new Rule()
                .pushDomain(strategy.proxyDomain)
                .pushIP(strategy.proxyIP)
            const direct = new Rule()
                .pushDomain(strategy.directDomain)
                .pushIP(strategy.directIP)
            return this._dns(hosts, proxy, direct, false)
        } else if (strategy.value >= 1000) { // 全部直接連接
            // 解析規則
            const proxy = new Rule()
                .pushDomain(strategy.proxyDomain)
                .pushIP(strategy.proxyIP)
            const direct = new Rule()
                .pushDomain(strategy.directDomain)
                .pushIP(strategy.directIP)
            return this._dns(hosts, proxy, direct, true)
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
            return this._dns(hosts, proxy, direct, false)
        }
        // 直連優先
        return this._dns(hosts, proxy, direct, true)
    }
    _dns(hosts, proxy, direct, usual) {
        const servers = []
        if (usual) {
            // 解析 西朝 域名
            if (direct.isValid()) {
                servers.push(...[
                    {
                        address: "119.29.29.29", // 騰訊
                        port: 53,
                        domains: direct.domain,
                        expectIPs: direct.ip,
                    },
                    {
                        address: "223.5.5.5", // 阿里
                        port: 53,
                        domains: direct.domain,
                        expectIPs: direct.ip,
                    },
                ])
            }
            // 解析 非西朝 域名
            if (proxy.isValid()) {
                servers.push(...[
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
                ])
            }
            // 未匹配的 使用西朝 dns
            servers.push(...[
                "119.29.29.29", // 騰訊
                "223.5.5.5", // 阿里
                "localhost",
            ])
        } else {
            // 解析 非西朝 域名
            if (proxy.isValid()) {
                servers.push(...[
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
                ])
            }
            // 解析 西朝 域名
            if (direct.isValid()) {
                servers.push(...[
                    {
                        address: "119.29.29.29", // 騰訊
                        port: 53,
                        domains: direct.domain,
                        expectIPs: direct.ip,
                    },
                    {
                        address: "223.5.5.5", // 阿里
                        port: 53,
                        domains: direct.domain,
                        expectIPs: direct.ip,
                    },
                ])
            }
            // 未匹配的 使用非西朝 dns
            servers.push(...[
                "8.8.8.8", // google
                "1.1.1.1", // cloudflare
                "https+local://doh.dns.sb/dns-query"
            ])
        }
        return {
            hosts: hosts,
            servers: servers
        }
    }
    inbounds() {
        const servers = []
        const ctx = this.ctx

        // 本地 socks5 代理
        const socks5 = ctx.socks5
        if (socks5) {
            const accounts = socks5.accounts
            servers.push({
                tag: "socks",
                protocol: "socks",
                listen: socks5.listen,
                port: socks5.port,
                settings: {
                    auth: accounts.length != 0 ? "password" : "noauth",
                    accounts: accounts.length != 0 ? accounts : undefined,
                    udp: socks5.udp,
                    userLevel: 0,
                }
            })
        }
        // 本地 http 代理
        const http = ctx.http
        if (http) {
            const accounts = http.accounts
            servers.push({
                tag: "http",
                protocol: "http",
                listen: http.listen,
                port: http.port,
                timeout: 300, // 超時時間爲 300 秒
                allowTransparent: false, // 爲 true 不止代理也轉發所有 http 請求
                accounts: accounts.length != 0 ? accounts : undefined,
                userLevel: 0,
            })
        }
        // 透明代理
        const proxy = ctx.proxy
        if (proxy) {
            servers.push({
                tag: "all-in",
                protocol: "dokodemo-door",
                listen: proxy.listen,
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
        return servers
    }
    outbounds() {
        const opts = this.opts
        const servers = []
        const direct = {
            tag: "direct",
            protocol: "freedom",
            streamSettings: {
                sockopt: {
                    mark: 2
                }
            }
        }
        const proxy = new OutboundReader(this.opts).toV2ray()
        if (opts.strategy.value < 900) {
            servers.push(proxy, direct)
        } else {
            servers.push(direct, proxy)
        }
        servers.push(
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
                        mark: 2
                    }
                }
            },
        )
        return servers
    }
    routing() {
        return {
            domainStrategy: "IPIfNonMatch",
            rules: this._rules(),
        }
    }
    _rules() {
        const strategy = this.opts.strategy
        if (strategy.value < 100) { // 全局代理
            // 解析規則
            const proxy = new Rule()
                .pushDomain(strategy.proxyDomain)
                .pushIP(strategy.proxyIP)
            const direct = new Rule()
                .pushDomain(strategy.directDomain)
                .pushIP(strategy.directIP)
            return this._createRules(proxy, direct, false)
        } else if (strategy.value < 200) {// 代理公有 ip
            // 解析規則
            const proxy = new Rule()
                .pushDomain(strategy.proxyDomain)
                .pushIP(strategy.proxyIP)
            const direct = new Rule()
                .pushIP([
                    "geoip:private",
                ])
                .pushDomain(strategy.directDomain)
                .pushIP(strategy.directIP)
            return this._createRules(proxy, direct, false)
        } else if (strategy.value >= 1000) { // 全部直接連接
            // 解析規則
            const proxy = new Rule()
                .pushDomain(strategy.proxyDomain)
                .pushIP(strategy.proxyIP)
            const direct = new Rule()
                .pushDomain(strategy.directDomain)
                .pushIP(strategy.directIP)
            return this._createRules(proxy, direct, true)
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
                "geoip:private",
            ])
            .pushDomain(strategy.directDomain)
            .pushIP(strategy.directIP)
        if (strategy.value < 900) { // 代理優先
            return this._createRules(proxy, direct, false)
        }
        // 直連優先
        return this._createRules(proxy, direct, true)
    }
    _createRules(proxy, direct, ok) {
        const rules = [
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
        ]
        // 處理 bt 協議
        const bittorrent = this.ctx.bittorrent
        if (bittorrent == "proxy" || bittorrent == "direct") {
            rules.push({
                type: "field",
                protocol: ["bittorrent"],
                outboundTag: bittorrent
            })
        }

        if (ok) {
            // 不代理 西朝
            if (direct.domain.length != 0) {
                rules.push({
                    type: "field",
                    domain: direct.domain,
                    outboundTag: "direct"
                })
            }
            if (direct.ip.length != 0) {
                rules.push({
                    type: "field",
                    ip: direct.ip,
                    outboundTag: "direct"
                })
            }
            // 代理 非西朝 域名
            if (proxy.domain.length != 0) {
                rules.push({
                    type: "field",
                    domain: proxy.domain,
                    outboundTag: "proxy"
                })
            }
            if (proxy.ip.length != 0) {
                rules.push({
                    type: "field",
                    ip: proxy.ip,
                    outboundTag: "proxy"
                })
            }
        } else {
            // 代理 非西朝 域名
            if (proxy.domain.length != 0) {
                rules.push({
                    type: "field",
                    domain: proxy.domain,
                    outboundTag: "proxy"
                })
            }
            if (proxy.ip.length != 0) {
                rules.push({
                    type: "field",
                    ip: proxy.ip,
                    outboundTag: "proxy"
                })
            }
            // 不代理 西朝
            if (direct.domain.length != 0) {
                rules.push({
                    type: "field",
                    domain: direct.domain,
                    outboundTag: "direct"
                })
            }
            if (direct.ip.length != 0) {
                rules.push({
                    type: "field",
                    ip: direct.ip,
                    outboundTag: "direct"
                })
            }
        }
        return rules
    }
}
class OutboundReader {
    constructor(opts) {
        if (opts instanceof ReaderOptions) {
            this.opts = opts
        } else {
            throw new Error("new OutboundReader(opts), must create from ReaderOptions")
        }
    }
    toV2ray() {
        const protocol = this.opts.outbound.protocol
        switch (protocol) {
            case "shadowsocks":
                return this.shadowsocks()
            case "trojan":
                return this.trojan()
            case "vless":
                return this.vless()
            case "vmess":
                return this.vmess()
            default:
                throw "not support protocol: " + protocol
        }
    }
    shadowsocks() {
        const opts = this.opts
        const outbound = opts.outbound
        return {
            tag: "proxy",
            protocol: outbound.protocol,
            settings: {
                servers: [
                    {
                        address: opts.ip,
                        port: outbound.port,
                        password: outbound.userID,
                        method: outbound.security,
                        uot: false,
                        level: outbound.level,
                    },
                ],
            },
            streamSettings: {
                sockopt: this._sockopt(),
            },
        }
    }
    trojan() {
        const opts = this.opts
        const outbound = opts.outbound
        const xtls = outbound.tls == "xtls"
        return {
            tag: "proxy",
            protocol: outbound.protocol,
            settings: {
                servers: [
                    {
                        address: opts.ip,
                        port: outbound.port,
                        password: outbound.userID,
                        flow: this._flow(),
                        level: outbound.level,
                    },
                ],
            },
            streamSettings: {
                network: "tcp",
                security: xtls ? "xtls" : "tls",
                tlsSettings: xtls ? undefined : this._tlsSettings(),
                xtlsSettings: xtls ? this._tlsSettings() : undefined,
                sockopt: this._sockopt(),
            },
        }
    }
    vless() {
        const opts = this.opts
        const outbound = opts.outbound
        return {
            tag: "proxy",
            protocol: outbound.protocol,
            settings: {
                vnext: [
                    {
                        address: opts.ip,
                        port: outbound.port,
                        users: [
                            {
                                id: outbound.userID,
                                flow: this._flow(),
                                encryption: "none",
                                level: outbound.level,
                            },
                        ],
                    },
                ],
            },
            streamSettings: {
                network: outbound.network,
                security: outbound.tls,
                tlsSettings: outbound.tls == "tls" ? this._tlsSettings() : undefined,
                xtlsSettings: outbound.tls == "xtls" ? this._tlsSettings() : undefined,
                tcpSettings: this._tcpSettings(),
                kcpSettings: this._kcpSettings(),
                wsSettings: this._wsSettings(),
                httpSettings: this._httpSettings(),
                dsSettings: this._dsSettings(),
                quicSettings: this._quicSettings(),
                sockopt: this._sockopt(),
            },
        }
    }
    vmess() {
        const opts = this.opts
        const outbound = opts.outbound
        return {
            tag: "proxy",
            protocol: outbound.protocol,
            settings: {
                vnext: [
                    {
                        address: opts.ip,
                        port: outbound.port,
                        users: [
                            {
                                id: outbound.userID,
                                alterId: outbound.alterID,
                                security: outbound.security,
                                level: outbound.level,
                            },
                        ],
                    },
                ],
            },
            streamSettings: {
                network: outbound.network,
                security: outbound.tls,
                tlsSettings: outbound.tls == "tls" ? this._tlsSettings() : undefined,
                xtlsSettings: outbound.tls == "xtls" ? this._tlsSettings() : undefined,
                tcpSettings: this._tcpSettings(),
                kcpSettings: this._kcpSettings(),
                wsSettings: this._wsSettings(),
                httpSettings: this._httpSettings(),
                dsSettings: this._dsSettings(),
                quicSettings: this._quicSettings(),
                sockopt: this._sockopt(),
            },
        }
    }
    _tlsSettings() {
        const outbound = this.opts.outbound
        const host = outbound.host ?? ''
        const address = outbound.address
        return {
            serverName: host == '' ? address : host,
            rejectUnknownSni: false,
            alpn: ["h2", "http/1.1"],
            // allowInsecure: true,//允許不安全的證書
            // "" | "chrome" | "firefox" | "safari" | "randomized"
            fingerprint: "firefox", // 模擬 tls 指紋
        }
    }
    _flow() {
        const outbound = this.opts.outbound
        if (outbound.tls != "xtls") {
            return
        } else if (outbound.protocol != "vless" && outbound.protocol != "trojan") {
            return
        }
        return outbound.flow
    }
    _sockopt() {
        return {
            mark: 2,
        }
    }


    _tcpSettings() {
        const opts = this.opts
        const outbound = opts.outbound
        if (outbound.network != "tcp") {
            return
        }
        return {
            header: {
                type: "none",
            },
        }
    }
    _kcpSettings() {
        const opts = this.opts
        const outbound = opts.outbound
        if (outbound.network != "kcp") {
            return
        }
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
    _wsSettings() {
        const opts = this.opts
        const outbound = opts.outbound
        if (outbound.network != "ws") {
            return
        }
        const path = outbound.path === "" ? "/" : outbound.path
        const host = outbound.host ?? ''
        const address = outbound.address
        return {
            path: path,
            headers: {
                Host: host == '' ? address : host,
            },
        }
    }
    _httpSettings() {
        const opts = this.opts
        const outbound = opts.outbound
        if (outbound.network != "http") {
            return
        }
        const path = outbound.path === "" ? "/" : outbound.path
        const host = outbound.host ?? ''
        const address = outbound.address
        return {
            path: path,
            method: "PUT",
            headers: {
                Host: host == '' ? address : host,
            },
        }
    }
    _dsSettings() {
        const opts = this.opts
        const outbound = opts.outbound
        if (outbound.network != "domainsocket") {
            return
        }
        const path = outbound.path === "" ? "/" : outbound.path
        return {
            path: path,
            abstract: false,
            padding: false,
        }
    }
    _quicSettings() {
        const opts = this.opts
        const outbound = opts.outbound
        if (outbound.network != "quic") {
            return
        }
        return {
            security: "none",
            key: "",
            header: {
                type: "none"
            }
        }
    }
}
