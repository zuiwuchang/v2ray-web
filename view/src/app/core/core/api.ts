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
    }
}
