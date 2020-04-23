import { isNumber, isObject, isString } from 'util'

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

    private _source = new Array<Element>()
    get source(): Array<Element> {
        return this._source
    }
    sort() {
        this._source.sort(Element.compare)
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
    static compare(l: Element, r: Element): number {
        if (l.outbound.name == r.outbound.name) {
            return l.id > r.id ? 1 : -1
        }
        return l.outbound.name > r.outbound.name ? 1 : -1
    }
}
export class Outbound {
    // 給人類看的 名稱
    name: string = ''

    // 連接地址
    add: string = ''
    // 連接端口
    port: string = ''
    // 連接主機名
    host: string = ''

    // 加密方案
    tls: string = ''

    // 使用的網路協議
    net: string = ''

    // websocket 請求路徑
    path: string = ''

    // 用戶身份識別碼
    userID: string = ''
    // 另外一個可選的用戶id
    alterID: string = ''
    // Security 加密方式
    security: string = ''
    // 用戶等級
    level: string = ''
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

        }
    }
}