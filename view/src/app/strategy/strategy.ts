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