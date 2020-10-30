import { RESTful } from './restful';
const root = '/api'
export const ServerAPI = {
    v1: {
        version: new RESTful(root, 'v1', 'version'),
        session: new RESTful(root, 'v1', 'session'),
        users: new RESTful(root, 'v1', 'users'),
        settings: new RESTful(root, 'v1', 'settings'),
        subscriptions: new RESTful(root, 'v1', 'subscriptions'),
        logs: new RESTful(root, 'v1', 'logs'),
        proxys: new RESTful(root, 'v1', 'proxys'),
    },
    static: {
        licenses: '/static/3rdpartylicenses.txt',
        license: '/static/LICENSE.txt',
    },
    v2ray: {
        settings: {
            get: `${root}/v2ray/settings/get`,
            put: `${root}/v2ray/settings/put`,
            test: `${root}/v2ray/settings/test`,
        },
    },
    proxy: {
        update: `${root}/proxy/update`,
        add: `${root}/proxy/add`,
        put: `${root}/proxy/put`,
        remove: `${root}/proxy/remove`,
        clear: `${root}/proxy/clear`,
        start: `${root}/proxy/start`,
        stop: `${root}/proxy/stop`,
        testOne: `${root}/proxy/test`,
        test: `${root}/ws/proxy/test`,
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