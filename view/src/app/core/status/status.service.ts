import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ServerAPI, getWebSocketAddr } from '../core/api';
import { isString, isNumber } from 'util';

export interface Status {
  run: boolean
  id?: number
  subscription?: number
  name?: string
}
@Injectable({
  providedIn: 'root'
})
export class StatusService {
  constructor() {
    this._do()
  }
  private _subject = new BehaviorSubject<Status>({
    run: false,
  })
  get observable(): Observable<Status> {
    return this._subject
  }
  private _websocket: WebSocket
  private _wait = 1
  private _timer
  private _do() {
    const addr = getWebSocketAddr(ServerAPI.proxy.status)
    console.info('ws connect', addr)

    const websocket = new WebSocket(addr)
    this._websocket = websocket
    websocket.onopen = (evt) => {
      if (this._websocket != websocket) {
        websocket.close()
        return
      }
      console.log("ws connect success", addr)
      this._wait = 1
      if (!this._timer) {
        this._timer = setInterval(() => {
          if (this._websocket) {
            this._websocket.send("ping")
          }
        }, 1000 * 10)
      }
    }
    websocket.onerror = (evt) => {
      this._onerror(evt, websocket)
    }
    websocket.onclose = (evt) => {
      this._onclose(evt, websocket)
    }
    websocket.onmessage = (evt) => {
      this._onmessage(evt, websocket)
    }
  }
  private _onerror(evt: Event, websocket: WebSocket) {
    websocket.close()
    if (this._websocket != websocket) {
      return
    }
    this._websocket = null
    if (this._timer) {
      clearInterval(this._timer)
      this._timer = null
    }
    this._rety()
  }
  private _onclose(evt: CloseEvent, websocket: WebSocket) {
    websocket.close()
    if (this._websocket != websocket) {
      return
    }
    this._websocket = null
    if (this._timer) {
      clearInterval(this._timer)
      this._timer = null
    }
    this._rety()
  }
  private _rety() {
    const wait = this._wait * 1000
    console.log("rety connect ws wait", wait)
    setTimeout(() => {
      this._do()
    }, wait)
    this._wait++
    if (this._wait > 5) {
      this._wait = 5
    }
  }
  private _onmessage(evt: MessageEvent, websocket: WebSocket) {
    if (this._websocket != websocket) {
      websocket.close()
      return
    }
    if (!isString(evt.data)) {
      return
    }
    try {
      const status = JSON.parse(evt.data) as Status
      this._notify(status)
    } catch (e) {
      console.warn(e)
    }
  }
  private _notify(status: Status) {
    const last = this._subject.value
    if (status.run) {
      if (last.run &&
        status.id == last.id &&
        status.subscription == last.subscription &&
        status.name == last.name) {
        return
      }
    } else if (!last.run) {
      return
    }
    if (status.run) {
      if (!isNumber(status.id)) {
        status.id = 0
      }
      if (!isNumber(status.subscription)) {
        status.subscription = 0
      }
    }
    this._subject.next({
      run: status.run,
      id: status.id,
      subscription: status.subscription,
      name: status.name,
    })
  }
}
