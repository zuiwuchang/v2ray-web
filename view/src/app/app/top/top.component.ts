import { Component, OnInit, OnDestroy } from '@angular/core';

import { StatusService, Status } from 'src/app/core/status/status.service';
import { Subscription } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { ServerAPI, getWebSocketAddr } from 'src/app/core/core/api';
import { isString } from 'king-node/dist/core';
import { SessionService } from 'src/app/core/session/session.service';
const MaxCount = 50

class Source {
  private _items = new Array<string>(MaxCount)
  private _index = 0
  private _count = 0
  private _text: string
  get text(): string {
    if (isString(this._text)) {
      return this._text
    }
    if (this._count == 0) {
      this._text = ''
    } else if (this._count == 1) {
      let index = this._index - 1
      if (index < 0) {
        index += this._items.length
      }
      this._text = this._items[index]
    } else {
      const arrs = new Array<string>(this._count)
      for (let i = 0; i < this._count; i++) {
        let index = this._index - 1 - i
        if (index < 0) {
          index += this._items.length
        }
        arrs[i] = this._items[index]
      }
      this._text = arrs.join("\n")
    }
    return this._text
  }
  push(v: string) {
    this._items[this._index] = v
    this._index++
    if (this._index == this._items.length) {
      this._index = 0
    }
    if (this._count < this._items.length) {
      this._count++
    }
    this._text = undefined
  }
}
@Component({
  selector: 'app-top',
  templateUrl: './top.component.html',
  styleUrls: ['./top.component.scss']
})
export class TopComponent implements OnInit, OnDestroy {
  constructor(private statusService: StatusService,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    public readonly sessionService: SessionService,
  ) { }
  private _closed = false
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  private _status: Status
  get status(): Status {
    return this._status
  }
  private _subscription: Subscription
  private _source = new Source()
  private _text: string = ''
  get text(): string {
    if (!this.pause) {
      this._text = this._source.text
    }
    return this._text
  }
  pause: boolean

  private _websocket: WebSocket
  private _wait = 1
  private _timer
  ngOnInit(): void {
    this._subscription = this.statusService.observable.subscribe((status) => {
      if (this._closed) {
        return
      }
      this._status = status
    })
    this._do()
  }
  ngOnDestroy() {
    this._closed = true
    this._subscription.unsubscribe()
    if (this._timer) {
      clearInterval(this._timer)
    }
  }
  onClickStop(evt: Event) {
    evt.stopPropagation()
    if (this._disabled) {
      return
    }
    this._disabled = true
    this.httpClient.get(ServerAPI.proxy.stop).toPromise().then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been stopped'),
      )
    }, (e) => {
      if (this._closed) {
        return
      }
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        e,
      )
    }).finally(() => {
      this._disabled = false
    })
  }
  private _do() {
    if (this._closed) {
      return
    }
    const addr = ServerAPI.v1.logs.websocketURL(null, {
      token: this.sessionService.token(),
    })
    console.info('logs ws connect', addr)

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
    const strs = evt.data.split("\n")
    for (let i = 0; i < strs.length; i++) {
      const str = strs[i].trim()
      if (str != "") {
        this._source.push(str)
      }
    }
  }
}
