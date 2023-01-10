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
        strategys: new RESTful(root, 'v1', 'strategys'),
        v2ray: new RESTful(root, 'v1', 'v2ray'),
        iptables: new RESTful(root, 'v1', 'iptables'),
    },
    static: {
        licenses: '/static/3rdpartylicenses.txt',
        license: '/static/LICENSE.txt',
    },
}
