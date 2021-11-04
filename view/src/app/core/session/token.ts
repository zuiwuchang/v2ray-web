import { isNumber, isString } from 'king-node/dist/core'
var Value: any
const TokenKey = 'token'
export class Token {
    constructor(public readonly value: string,
        public readonly at: number,
        public readonly maxage: number,
    ) {
    }
    isExpired(): boolean {
        const expired = this.at + this.maxage * 1000
        return expired <= new Date().getTime()
    }
}

export function SaveToken(value: string, at: number, maxage: number) {
    if (typeof localStorage == 'undefined') {
        return
    }
    const token = LoadToken()
    if (token && token.at >= at) {
        return
    }
    try {
        localStorage.setItem(TokenKey, JSON.stringify({
            value: value,
            at: at,
            maxage: maxage,
        }))
    } catch (e) {
        console.log(`save token error : `, e)
    }
}
export function LoadToken(): Token {
    try {
        if (typeof localStorage == 'undefined') {
            return Value
        }
        const str = localStorage.getItem(TokenKey)
        if (!isString(str)) {
            return Value
        }
        const obj = JSON.parse(str as string)
        if (obj && isNumber(obj.at) && isNumber(obj.maxage) && isString(obj.value)) {
            const token = new Token(obj.value, obj.at, obj.maxage)
            if (!token.isExpired()) {
                return token
            }
        }
    } catch (e) {
        console.log(`load token error : `, e)
    }
    return Value
}
export function DeleteToken(value: string) {
    try {
        if (typeof localStorage == 'undefined') {
            return
        }
        const token = LoadToken()
        if (token && token.value == value) {
            localStorage.removeItem(TokenKey)
        }
    } catch (e) {
        console.log(`delete token error : `, e)
    }
}
