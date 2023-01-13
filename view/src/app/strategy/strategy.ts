export interface Strategy {
    // 唯一的名稱 供人類查看
    // 名稱 default 是系統保留的策略，其它策略將繼承這個策略 Name 和 Value 之外的所有值
    name: string

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
    value: number

    // 靜態 ip 列表
    // baidu.com 127.0.0.1
    // dns.google 8.8.8.8 8.8.4.4
    host: string

    // 這些 ip 使用代理
    proxyIP: string
    // 這些 域名 使用代理
    proxyDomain: string

    // 這些 ip 直連
    directIP: string
    // 這些 域名 直連
    directDomain: string

    // 這些 ip 阻止訪問
    blockIP: string
    // 這些 域名 阻止訪問
    blockDomain: string
}
export class StrategyValue {
    static fromStrategy(o: Strategy): StrategyValue {
        const v = new StrategyValue()
        v.name = o.name
        v.value = o.value
        v.host = o.host
        v.proxy.set(o.proxyIP, o.proxyDomain)
        v.direct.set(o.directIP, o.directDomain)
        v.block.set(o.blockIP, o.blockDomain)
        return v
    }
    name: string
    value: number
    host: string
    proxy: StrategyElement
    direct: StrategyElement
    block: StrategyElement
    constructor(opts?: {
        name?: string
        value?: number
        host?: string
        proxy?: StrategyElement
        direct?: StrategyElement
        block?: StrategyElement
    }) {
        this.name = opts?.name ?? ''
        this.value = opts?.value ?? 0
        this.host = opts?.host ?? ''
        this.proxy = opts?.proxy ?? new StrategyElement()
        this.direct = opts?.direct ?? new StrategyElement()
        this.block = opts?.block ?? new StrategyElement()
    }
    cloneTo(o: StrategyValue) {
        o.name = this.name
        o.value = this.value
        o.host = this.host
        this.proxy.cloneTo(o.proxy)
        this.direct.cloneTo(o.direct)
        this.block.cloneTo(o.block)
    }
}
export class StrategyElement {
    ip: string
    domain: string
    constructor(opts?: {
        ip?: string
        domain?: string
    }) {
        this.ip = opts?.ip ?? ''
        this.domain = opts?.domain ?? ''
    }
    set(ip: string, domain: string) {
        this.ip = ip
        this.domain = domain
    }
    cloneTo(o: StrategyElement) {
        o.ip = this.ip
        o.domain = this.domain
    }
}