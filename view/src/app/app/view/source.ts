import { isBoolean, isNumber, isObject, isString } from 'util'
import { Base64 } from 'js-base64';
export class Source {
    private _items = new Array<Panel>()
    private _keys = new Map<number, Panel>()
    get items(): Array<Panel> {
        return this._items
    }
    put(panel: Panel) {
        if (this._keys.has(panel.id)) {
            console.warn(`${panel.id} panel already exists`)
            return
        }
        this._keys.set(panel.id, panel)
        this._items.push(panel)
    }
    set(element: Element) {
        const panel = this._getPanel(element.subscription)
        panel.source.push(element)
    }
    private _getPanel(id: number): Panel {
        let panel = this._keys.get(id)
        if (!panel) {
            panel = this._keys.get(0)
        }
        return panel
    }
    sort() {
        const items = this._items
        items.sort(Panel.compare)
        for (let i = 0; i < items.length; i++) {
            items[i].sort()
        }
    }
}
export class Panel {
    id: number
    name: string
    source = new Array<Element>()
    sort() {
        this.source.sort(Element.compare)
    }
    static compare(l: Panel, r: Panel): number {
        if (l.id == 0) {
            return -1
        } else if (r.id == 0) {
            return 1
        }
        if (l.name == r.name) {
            return l.id > r.id ? 1 : -1
        }
        return l.name > r.name ? 1 : -1
    }
}
export class Element {
    id: number = 0
    subscription: number = 0
    outbound: Outbound
    // 是否 正在發送請求
    request?: boolean
    duration?: number
    error?: string

    constructor(net?: Element) {
        if (isObject(net)) {
            if (isNumber(net.id)) {
                this.id = net.id
            }
            if (isNumber(net.subscription)) {
                this.subscription = net.subscription
            }
            this.outbound = new Outbound(net.outbound)
        } else {
            this.outbound = new Outbound()
        }
    }
    static compareDuration(l: Element, r: Element): number {
        let ld = l.duration
        let rd = r.duration
        if (!isNumber(ld) || ld <= 0) {
            ld = 1000 * 60 * 60
        }
        if (!isNumber(rd) || rd <= 0) {
            rd = 1000 * 60 * 60
        }
        if (ld == rd) {
            return Element.compare(l, r)
        }
        return ld > rd ? 1 : -1
    }
    static compare(l: Element, r: Element): number {
        if (l.outbound.name != r.outbound.name) {
            return l.outbound.name > r.outbound.name ? 1 : -1
        }
        if (l.outbound.net != r.outbound.net) {
            return l.outbound.net > r.outbound.net ? 1 : -1
        }
        if (l.outbound.tls != r.outbound.tls) {
            return l.outbound.tls > r.outbound.tls ? 1 : -1
        }
        if (l.id != r.id) {
            return l.id > r.id ? 1 : -1
        }
        return 0
    }
}
export class Outbound {
    // 給人類看的 名稱
    name: string = ''

    // 連接地址
    add: string = ''
    // 連接端口
    port: string = '80'
    // 連接主機名
    host: string = ''

    // 加密方案
    tls: string = 'none'

    // 使用的網路協議
    net: string = 'ws'

    // websocket 請求路徑
    path: string = '/'

    // 用戶身份識別碼
    userID: string = 'xxxxxxxx'
    // 另外一個可選的用戶id
    alterID: string = '0'
    // Security 加密方式
    security: string = 'auto'
    // 用戶等級
    level: string = '0'

    vless: boolean = false
    constructor(net?: Outbound) {
        if (isObject(net)) {
            if (isString(net.name)) {
                this.name = net.name
            }
            if (isString(net.add)) {
                this.add = net.add
            }
            if (isString(net.port)) {
                this.port = net.port
            }
            if (isString(net.host)) {
                this.host = net.host
            }
            if (isString(net.tls)) {
                this.tls = net.tls
            }
            if (isString(net.net)) {
                this.net = net.net
            }
            if (isString(net.path)) {
                this.path = net.path
            }
            if (isString(net.userID)) {
                this.userID = net.userID
            }
            if (isString(net.alterID)) {
                this.alterID = net.alterID
            }
            if (isString(net.security)) {
                this.security = net.security
            }
            if (isString(net.level)) {
                this.level = net.level
            }
            if (isBoolean(net.vless)) {
                this.vless = net.vless
            }
        }
    }
    format() {
        this.port = this._valNumber(this.port)
        this.alterID = this._valNumber(this.alterID)
        this.level = this._valNumber(this.level)
    }
    private _valNumber(v: any): string {
        if (isNumber(v)) {
            return v.toString()
        }
        return v
    }
    cloneTo(other: Outbound) {
        other.name = this.name
        other.add = this.add
        other.port = this.port
        other.host = this.host
        other.tls = this.tls
        other.net = this.net
        other.path = this.path
        other.userID = this.userID
        other.alterID = this.alterID
        other.security = this.security
        other.level = this.level
        other.vless = this.vless
    }
    equal(other: Outbound): boolean {
        return other.name == this.name &&
            other.add == this.add &&
            other.port == this.port &&
            other.host == this.host &&
            other.tls == this.tls &&
            other.net == this.net &&
            other.path == this.path &&
            other.userID == this.userID &&
            other.alterID == this.alterID &&
            other.security == this.security &&
            other.level == this.level &&
            other.vless == this.vless
    }

    static fromBase64(str: string): Outbound {
        str = str.replace(/=+/, '')
        str = Base64.decode(str)
        const obj = JSON.parse(str)
        const outbound = new Outbound()

        outbound.name = obj.ps
        outbound.add = obj.add
        outbound.port = obj.port
        outbound.host = obj.host
        outbound.tls = obj.tls
        outbound.net = obj.net
        outbound.path = obj.path
        outbound.userID = obj.id
        outbound.alterID = obj.aid
        outbound.security = obj.type
        outbound.level = obj.v
        outbound.vless = obj.vless
        return outbound
    }
    toBase64(): string {
        const str = JSON.stringify({
            ps: this.name,
            add: this.add,
            port: this.port,
            host: this.host,
            tls: this.tls,
            net: this.net,
            path: this.path,
            id: this.userID,
            aid: this.alterID,
            type: this.security,
            v: this.level,
            vless: this.vless ? true : false,
        })
        return Base64.encode(str)
    }
}