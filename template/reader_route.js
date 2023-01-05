function renderRouting(ctx) {
    return {
        domainStrategy: "AsIs",
        domainMatcher: "hybrid",
        rules: [
            // 私有ip 和 朝鮮ip 不使用代理
            {
                type: "field",
                ip: ["geoip:private", "geoip:cn", "223.5.5.5"],
                outboundTag: "freedom",
            },
            // 非朝鮮域名使用代理
            {
                type: "field",
                domain: [
                    "geosite:geolocation-!cn",
                ],
                outboundTag: "proxy",
            },
            // 屏蔽廣告
            {
                type: "field",
                domain: ["geosite:category-ads-all"],
                outboundTag: "blackhole",
            },
            // 朝鮮不使用代理
            {
                type: "field",
                domain: ["geosite:cn"],
                outboundTag: "freedom",
            },
            // bt 下載不使用代理
            {
                type: "field",
                protocol: ["bittorrent"],
                outboundTag: "freedom",
            },

        ],

    }
}