const root = '/api'
export const ServerAPI = {
    version: `${root}/app/version`,
    login: `${root}/app/login`,
    restore: `${root}/app/restore`,
    logout: `${root}/app/logout`,
    logs: `${root}/ws/app/logs`,
    settings: {
        get: `${root}/settings/get`,
        put: `${root}/settings/put`,
    },
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
        testOne: `${root}/proxy/test`,
        test: `${root}/ws/proxy/test`,
        status: `${root}/ws/proxy/status`,
    },
    iptables: {
        view: `${root}/iptables/view`,
        get: `${root}/iptables/get`,
        getDefault: `${root}/iptables/get/default`,
        put: `${root}/iptables/put`,
        init: `${root}/iptables/init`,
        restore: `${root}/iptables/restore`,
    },
}
export function getWebSocketAddr(path: string): string {
    const location = document.location
    let addr: string
    if (location.protocol == "https") {
        addr = `wss://${location.hostname}`
        if (location.port == "") {
            addr += ":443"
        } else {
            addr += `:${location.port}`
        }
    } else {
        addr = `ws://${location.hostname}`
        if (location.port == "") {
            addr += ":80"
        } else {
            addr += `:${location.port}`
        }
    }
    return `${addr}${path}`
}