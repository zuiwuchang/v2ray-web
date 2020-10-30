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
        v2ray: new RESTful(root, 'v1', 'v2ray'),
        iptables: new RESTful(root, 'v1', 'iptables'),
    },
    static: {
        licenses: '/static/3rdpartylicenses.txt',
        license: '/static/LICENSE.txt',
    },
    proxy: {
        update: `${root}/proxy/update`,
        add: `${root}/proxy/add`,
        put: `${root}/proxy/put`,
        clear: `${root}/proxy/clear`,
        start: `${root}/proxy/start`,
        stop: `${root}/proxy/stop`,
        testOne: `${root}/proxy/test`,
        test: `${root}/ws/proxy/test`,
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