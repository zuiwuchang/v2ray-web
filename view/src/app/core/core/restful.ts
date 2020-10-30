import { HttpHeaders, HttpParams, HttpClient } from '@angular/common/http'
import { isNumber, isString, isObject, isArray, isUndefinedOrNull } from "king-node/dist/core";
export function resolveError(e): string {
    if (!e) {
        return "nil"
    }
    if (isString(e)) {
        return e
    }
    if (isObject(e) && isNumber(e.status)) {
        return resolveHttpError(e)
    }
    return "unknow"
}
export function resolveHttpError(e) {
    if (isString(e.error)) {
        return `${e.status} ${e.error}`
    }
    if (e.error) {
        if (e.error.description) {
            return `${e.status} ${e.error.description}`
        } else if (e.error.error) {
            return `${e.status} ${e.error.error}`
        }
    } else if (e.message) {
        return `${e.status} ${e.message}`
    }
    return `${e.status} ${e.statusText}`
}
export function wrapPromise<T>(promise: Promise<T>): Promise<T> {
    return new Promise<T>(function (resolve, reject) {
        promise.then(
            function (v) {
                resolve(v)
            },
            function (e) {
                reject(resolveError(e))
            },
        )
    })
}
export class RESTful {
    constructor(public root, public version, public url: string) {

    }
    get baseURL(): string {
        return `${this.root}/${this.version}/${this.url}`
    }
    oneURL(id: string | number | boolean | Array<any>): string {
        let val: string
        if (isArray(id)) {
            val = (id as Array<any>).map<string>((val) => encodeURIComponent(encodeURIComponent(val))).join('/')
        } else {
            val = encodeURIComponent(encodeURIComponent(id as string))
        }
        return `${this.baseURL}/${val}`
    }
    onePatchURL(id: string | number | boolean | Array<any>, patch: string): string {
        return `${this.oneURL(id)}/${patch}`
    }
    websocketURL(id: string | number | boolean | Array<any>, params?: {
        [param: string]: string;
    }): string {
        const location = document.location
        let addr: string
        // console.log(location.protocol)
        if (location.protocol == "https:") {
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
        let val: string
        if (!isUndefinedOrNull(id)) {
            if (isArray(id)) {
                val = (id as Array<any>).map<string>((val) => encodeURIComponent(encodeURIComponent(val))).join('/')
            } else {
                val = encodeURIComponent(encodeURIComponent(id as string))
            }
        }
        let url = `${addr}${this.baseURL}`
        if (!isUndefinedOrNull(val)) {
            url += '/' + val
        }
        url += '/websocket'
        if (params) {
            const query = new URLSearchParams('')
            for (const key in params) {
                query.set(key, params[key])
            }
            return `${url}?${query.toString()}`
        }
        return url
    }
    get<T>(client: HttpClient, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T>;
    get(client: HttpClient, options: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType: 'text';
        withCredentials?: boolean;
    }): Promise<string>;
    get(client: HttpClient, options?: any): any {
        return wrapPromise(client.get(this.baseURL, options).toPromise())
    }
    post<T>(client: HttpClient, body: any | null, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T> {
        return wrapPromise(client.post<T>(this.baseURL, body, options).toPromise())
    }
    delete<T>(client: HttpClient, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T> {
        return wrapPromise(client.delete<T>(this.baseURL, options).toPromise())
    }
    put<T>(client: HttpClient, body: any | null, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T> {
        return wrapPromise(client.put<T>(this.baseURL, body, options).toPromise())
    }
    patch<T>(client: HttpClient, patch: string, body: any | null, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T> {
        return wrapPromise(client.patch<T>(`${this.baseURL}/${patch}`, body, options).toPromise())
    }
    getOne<T>(client: HttpClient, id: string | number | boolean | Array<any>, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T>;
    getOne(client: HttpClient, id: string | number | boolean | Array<any>, options: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType: 'text';
        withCredentials?: boolean;
    }): Promise<string>;
    getOne(client: HttpClient, id: string | number | boolean | Array<any>, options?: any): any {
        return wrapPromise(client.get(this.oneURL(id), options).toPromise())
    }
    postOne<T>(client: HttpClient, id: string | number | boolean | Array<any>, body: any | null, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T> {
        return wrapPromise(client.post<T>(this.oneURL(id), body, options).toPromise())
    }
    deleteOne<T>(client: HttpClient, id: string | number | boolean | Array<any>, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T> {
        return wrapPromise(client.delete<T>(this.oneURL(id), options).toPromise())
    }
    putOne<T>(client: HttpClient, id: string | number | boolean | Array<any>, body: any | null, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T> {
        return wrapPromise(client.put<T>(this.oneURL(id), body, options).toPromise())
    }
    patchOne<T>(client: HttpClient, id: string | number | boolean | Array<any>, patch: string, body: any | null, options?: {
        headers?: HttpHeaders | {
            [header: string]: string | string[];
        };
        observe?: 'body';
        params?: HttpParams | {
            [param: string]: string | string[];
        };
        reportProgress?: boolean;
        responseType?: 'json';
        withCredentials?: boolean;
    }): Promise<T> {
        return wrapPromise(client.patch<T>(this.onePatchURL(id, patch), body, options).toPromise())
    }
}
