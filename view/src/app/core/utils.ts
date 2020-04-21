import { isString, isNumber } from 'util';

export namespace Utils {
    export function resolveError(e): string {
        if (!e) {
            return "nil"
        }
        if (isNumber(e.status)) {
            return resolveHttpError(e)
        }
        return "unknow"
    }
    export function resolveHttpError(e) {
        if (e.status != 500) {
            return `${e.status} ${e.statusText}`
        }
        if (isString(e.error)) {
            return `${e.status} ${e.error}`
        }
        if (e.error) {
            return `${e.status} ${e.error.description}`
        } else if (e.message) {
            return `${e.status} ${e.message}`
        }
        return `${e.status}`
    }
}
