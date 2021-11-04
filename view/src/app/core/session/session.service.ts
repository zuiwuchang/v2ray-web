import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from '../core/api';
import { BehaviorSubject, Observable } from 'rxjs';
import { Completer } from 'king-node/dist/async/completer';
import { Mutex } from 'king-node/dist/async/sync';
import { isString, Exception, isNumber } from 'king-node/dist/core';
import { SaveToken, LoadToken, DeleteToken } from './token'
var Value: any
interface Response {
  session: Session
  token: string
  maxage: number
}
export class Session {
  name: string = ''
  root = false
  token: string = ''
}

@Injectable({
  providedIn: 'root'
})
export class SessionService {
  constructor(private httpClient: HttpClient) {
    // 恢復 session
    this._restore();
  }
  private readonly mutex_ = new Mutex()
  private readonly subject_ = new BehaviorSubject<Session>(Value)
  private readonly ready_ = new Completer<boolean>()
  get observable(): Observable<Session> {
    return this.subject_
  }
  get ready(): Promise<boolean> {
    return this.ready_.promise
  }
  private async _restore() {
    console.log('start session restore')
    await this.mutex_.lock()
    const token = LoadToken()
    try {
      if (token) {
        const response = await ServerAPI.v1.session.getOne<Session>(this.httpClient, [Math.floor(token.at / 1000), token.maxage, token.value])
        if (response && isString(response.name)) {
          console.info(`session restore`, response)
          response.token = token.value
          this.subject_.next(response)
        }
      }
    } catch (e) {
      console.error(`restore error : `, e)
    } finally {
      this.mutex_.unlock()
      this.ready_.resolve(true)
    }
  }
  /**
 * 登入
 * @param name 
 * @param password 
 * @param remember 
 */
  async login(name: string, password: string, remember: boolean): Promise<Session> {
    await this.mutex_.lock()
    let result: any// Session
    try {
      const timestamp = new Date().getTime()
      const response = await ServerAPI.v1.session.post<Response>(this.httpClient, {
        name: name,
        password: password,
      })
      if (response) {
        if (remember && isString(response.token) && isNumber(response.maxage)) {
          SaveToken(response.token, timestamp, response.maxage)
        }
        response.session.token = response.token
        this.subject_.next(response.session)
      } else {
        console.warn(`login unknow result`, response)
        throw new Exception("login unknow result")
      }
    } finally {
      this.mutex_.unlock()
    }
    return result
  }


  /**
  * 登出
  */
  async logout() {
    await this.mutex_.lock()
    try {
      if (this.subject_.value == null) {
        return
      }
      DeleteToken(this.subject_.value.token)
      this.subject_.next(Value)
    } finally {
      this.mutex_.unlock()
    }
  }
  token(): string {
    if (this.subject_.value) {
      return this.subject_.value.token
    }
    return Value
  }
}
