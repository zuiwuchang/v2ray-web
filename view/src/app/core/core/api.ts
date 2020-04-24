const root = '/api'
export const ServerAPI = {
    login: `${root}/app/login`,
    restore: `${root}/app/restore`,
    logout: `${root}/app/logout`,

    user: {
        list: `${root}/user/list`,
        add: `${root}/user/add`,
        remove: `${root}/user/remove`,
        password: `${root}/user/password`,
    },
    v2ray: {
        settings: {
            get: `${root}/v2ray/settings/get`,
            put: `${root}/v2ray/settings/put`,
            test: `${root}/v2ray/settings/test`,
        },
        subscription: {
            list: `${root}/v2ray/subscription/list`,
            put: `${root}/v2ray/subscription/put`,
            add: `${root}/v2ray/subscription/add`,
            remove: `${root}/v2ray/subscription/remove`,
        },
    },
    proxy: {
        list: `${root}/proxy/list`,
        update: `${root}/proxy/update`,
        add: `${root}/proxy/add`,
        put: `${root}/proxy/put`,
        remove: `${root}/proxy/remove`,
        clear: `${root}/proxy/clear`,
        start: `${root}/proxy/start`,
        stop: `${root}/proxy/stop`,
        test: `${root}/ws/proxy/test`,
        status: `${root}/ws/proxy/status`,
    },
}
export function getWebSocketAddr(path: string): string {
    const location = document.location
    let addr: string
    if (location.protocol == "https") {
        addr = `wss://${location.host}:${location.port}`
        if (location.port == "") {
            addr += "443"
        }
    } else {
        addr = `ws://${location.host}:${location.port}`
        if (location.port == "") {
            addr += "80"
        }
    }
    return `${addr}${path}`
}