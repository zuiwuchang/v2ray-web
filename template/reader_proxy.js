function renderRouting(ctx) {
    return {
        domainStrategy: "AsIs",
        domainMatcher: "hybrid",
        rules: [
            {
                type: "field",
                outboundTag: "proxy",
            },
        ],
    }
}