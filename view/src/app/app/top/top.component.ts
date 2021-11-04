import { Component, OnInit, OnDestroy, AfterViewInit, ViewChild, ElementRef } from '@angular/core';

import { StatusService, Status } from 'src/app/core/status/status.service';
import { fromEvent, Subscription } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { ServerAPI } from 'src/app/core/core/api';
import { isString } from 'king-node/dist/core';
import { SessionService } from 'src/app/core/session/session.service';
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { takeUntil } from 'rxjs/operators';
import { Closed } from 'src/app/core/core/utils';
var Value: any
@Component({
  selector: 'app-top',
  templateUrl: './top.component.html',
  styleUrls: ['./top.component.scss']
})
export class TopComponent implements OnInit, OnDestroy, AfterViewInit {
  constructor(private statusService: StatusService,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    public readonly sessionService: SessionService,
  ) { }
  private _closed = new Closed()
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  private _status: Status = { run: false }
  get status(): Status {
    return this._status
  }
  private _subscription: Subscription = Value

  private _websocket: WebSocket = Value
  private _wait = 1
  private _timer = Value
  private _xterm: Terminal = Value
  private _fitAddon = new FitAddon()
  private _webLinksAddon = new WebLinksAddon()
  ngOnInit(): void {
    this._subscription = this.statusService.observable.pipe(
      takeUntil(this._closed.observable)
    ).subscribe((status) => {
      this._status = status
    })
    this._do()
  }
  ngOnDestroy() {
    this._closed.close()
    this._subscription.unsubscribe()
    if (this._timer) {
      clearInterval(this._timer)
    }
    this._xterm.dispose()
    this._fitAddon.dispose()
    this._webLinksAddon.dispose()
  }
  @ViewChild("xterm")
  xterm: ElementRef = Value
  ngAfterViewInit() {
    const xterm = new Terminal({
      cursorBlink: false,
      screenReaderMode: true,
    })
    this._xterm = xterm
    xterm.loadAddon(this._fitAddon)
    xterm.loadAddon(this._webLinksAddon)
    xterm.open(this.xterm.nativeElement)
    this._fitAddon.fit()
    this._do()
    fromEvent(window, 'resize').pipe(
      takeUntil(this._closed.observable),
    ).subscribe(() => {
      this._fitAddon.fit()
    })
  }
  onClickStop(evt: Event) {
    evt.stopPropagation()
    if (this._disabled) {
      return
    }
    this._disabled = true
    ServerAPI.v1.proxys.postOne(this.httpClient, 'stop', null).then(() => {
      if (this._closed.isClosed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been stopped'),
      )
    }, (e) => {
      if (this._closed.isClosed) {
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
    if (this._closed.isClosed) {
      return
    }
    const addr = ServerAPI.v1.logs.websocketURL(Value, {
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
    this._websocket = Value
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
    this._websocket = Value
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
        this._xterm.writeln(str)
      }
    }
  }
}
