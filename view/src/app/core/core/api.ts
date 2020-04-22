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
            list: `${root}/v2ray/subscription/list1`,
            get: `${root}/v2ray/subscription/get1`,
            put: `${root}/v2ray/subscription/put1`,
            add: `${root}/v2ray/subscription/add1`,
        },
    },
}
