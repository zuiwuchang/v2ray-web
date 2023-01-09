import { Component, OnInit, Input, OnDestroy, ViewChild, ElementRef, AfterViewInit } from '@angular/core';
import { Panel, Element } from '../view/source';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { isArray, isString, isNumber } from 'king-node/dist/core';
import { MatDialog } from '@angular/material/dialog';
import { AddComponent } from '../add/add.component';
import { EditComponent } from '../edit/edit.component';
import { ConfirmComponent } from 'src/app/shared/dialog/confirm/confirm.component';
import { StatusService, Status } from 'src/app/core/status/status.service';
import { Subscription } from 'rxjs';
import * as ClipboardJS from 'clipboard'
import { QrcodeComponent } from '../dialog/qrcode/qrcode.component';
import { SessionService } from 'src/app/core/session/session.service';
// 正在運行
const StatusRunning = 1
// 錯誤
const StatusError = 2
// 完成
const StatusOk = 3
interface Message {
  status: number
  id: number
  error?: string
  duration?: number
}
var Value: any
@Component({
  selector: 'app-view-panel',
  templateUrl: './view-panel.component.html',
  styleUrls: ['./view-panel.component.scss']
})
export class ViewPanelComponent implements OnInit, OnDestroy, AfterViewInit {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialog: MatDialog,
    private statusService: StatusService,
    public readonly sessionService: SessionService,
  ) { }
  @Input('panel')
  panel: Panel = Value
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  private _closed = false
  private _websocket: WebSocket = Value
  private _status: Status = Value
  get status(): Status {
    return this._status
  }
  private _subscription: Subscription = Value
  ngOnInit(): void {
    this._subscription = this.statusService.observable.subscribe((status) => {
      if (this._closed) {
        return
      }
      this._status = status
    })
  }
  get isSubscription(): boolean {
    if (this._status.run) {
      return this._status.subscription == this.panel.id
    }
    return false
  }
  getIcon(element: Element): string {
    if (this.isCurrent(element)) {
      return 'done'
    }
    if (element && element.outbound) {
      if (element.outbound.protocol == "vless") {
        return 'star_half'
      } else if (element.outbound.protocol == "vmess") {
        return 'star'
      } else if (element.outbound.protocol == "shadowsocks") {
        return 'flight'
      } else if (element.outbound.protocol == "trojan") {
        return 'https'
      }
    }
    return 'feedback'
  }
  isCurrent(element: Element) {
    if (this.status.run) {
      return this._status.subscription == this.panel.id && this._status.id == element.id
    }
    return false
  }
  ngOnDestroy() {
    this._closed = true
    this._subscription.unsubscribe()
    if (this._websocket) {
      const websocket = this._websocket
      this._websocket = Value
      websocket.close()
    }
    const source = this.panel.source
    for (let i = 0; i < source.length; i++) {
      source[i].request = undefined
    }
    if (this._clipboard) {
      this._clipboard.destroy();
    }
  }
  private _clipboard: any = null
  @ViewChild("btnClipboard")
  private _btnClipboard: ElementRef = Value
  ngAfterViewInit() {
    this._clipboard = new ClipboardJS(this._btnClipboard.nativeElement).on('success', () => {
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get("data copied"),
      )
    }).on('error', (evt) => {
      console.error('Action:', evt.action);
      console.error('Trigger:', evt.trigger);
    });
  }
  onClickSort() {
    this.panel.source.sort(Element.compareDuration)
  }
  onClickTest() {
    this._disabled = true
    if (this._websocket) {
      this._websocket.close()
    }

    const addr = ServerAPI.v1.proxys.websocketURL('test', {
      token: this.sessionService.token(),
    })
    const websocket = new WebSocket(addr)
    this._websocket = websocket
    websocket.onerror = (evt) => {
      this._onerror(evt, websocket)
    }
    websocket.onclose = (evt) => {
      this._onclose(evt, websocket)
    }
    websocket.onopen = () => {
      this._onopen(websocket)
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
    console.warn(evt)
    this.toasterService.pop('error',
      this.i18nService.get('error'),
      'websocket error',
    )
    this._resetWebsocket()
  }
  private _onclose(evt: CloseEvent, websocket: WebSocket) {
    websocket.close()
    if (this._websocket != websocket) {
      return
    }
    console.warn(evt)
    this.toasterService.pop('error',
      this.i18nService.get('error'),
      'websocket closed',
    )
    this._resetWebsocket()
  }
  private _onopen(websocket: WebSocket) {
    if (this._websocket != websocket) {
      websocket.close()
      return
    }
    const source = this.panel.source
    const items = new Array<string>()
    for (let i = 0; i < source.length; i++) {
      const element = source[i]
      try {
        const str = JSON.stringify({
          id: element.id,
          outbound: element.outbound,
        })
        items.push(str)
        element.request = undefined
        element.error = undefined
        element.duration = undefined
      } catch (e) {
        console.warn(e)
      }
    }
    try {
      for (let i = 0; i < items.length; i++) {
        websocket.send(items[i])
      }
      websocket.send("close")
    } catch (e) {
      console.warn(e)
      if (websocket == this._websocket) {
        this._resetWebsocket()
      }
      websocket.close()
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
    if (evt.data == "close") {
      websocket.close()
      this._resetWebsocket()
      return
    }

    try {
      const resopnse = JSON.parse(evt.data) as Message
      switch (resopnse.status) {
        case StatusOk:
          this._setOk(resopnse)
          break
        case StatusError:
          this._setError(resopnse)
          break
        case StatusRunning:
          this._setRunning(resopnse)
          break
      }
    } catch (e) {
      console.log(e)
    }
  }
  private _setOk(msg: Message) {
    const source = this.panel.source
    for (let i = 0; i < source.length; i++) {
      if (source[i].id == msg.id) {
        source[i].request = undefined
        source[i].error = undefined
        source[i].duration = msg.duration
        break
      }
    }
  }
  private _setError(msg: Message) {
    const source = this.panel.source
    for (let i = 0; i < source.length; i++) {
      if (source[i].id == msg.id) {
        source[i].request = undefined
        source[i].error = msg.error
        source[i].duration = undefined
        break
      }
    }
  }
  private _setRunning(msg: Message) {
    const source = this.panel.source
    for (let i = 0; i < source.length; i++) {
      if (source[i].id == msg.id) {
        source[i].request = true
        break
      }
    }
  }
  private _resetWebsocket() {
    const source = this.panel.source
    for (let i = 0; i < source.length; i++) {
      source[i].request = undefined
    }
    this._disabled = false
    this._websocket = Value
  }
  onClickAdd() {
    this.matDialog.open(AddComponent, {
      data: this.panel,
      disableClose: true,
    })
  }
  onClickClear() {
    this.matDialog.open(ConfirmComponent, {
      data: {
        title: this.i18nService.get("clear proxy element title"),
        content: this.i18nService.get("clear proxy element"),
      },
    }).afterClosed().toPromise().then((data) => {
      if (this._closed || !data) {
        return
      }
      this._clear()
    })
  }
  private _clear() {
    this._disabled = true
    ServerAPI.v1.proxys.deleteOne(this.httpClient, 'clear', {
      params: {
        subscription: this.panel.id.toString(),
      }
    }).then(() => {
      if (this._closed) {
        return
      }
      this.panel.source = new Array<Element>()
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been cleared'),
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
  onClickUpdate() {
    this._disabled = true
    ServerAPI.v1.proxys.postOne<Array<any>>(this.httpClient, 'update', {
      id: this.panel.id,
    }).then((data) => {
      if (this._closed) {
        return
      }
      const source = new Array<Element>()
      this.panel.source.slice(0, this.panel.source.length)
      if (isArray(data) && data.length > 0) {
        for (let i = 0; i < data.length; i++) {
          source.push(new Element(data[i]))
        }
        source.sort(Element.compare)
      }
      this.panel.source = source
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been updated'),
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
  onClickStart(element: Element) {
    this._disabled = true
    ServerAPI.v1.proxys.postOne(this.httpClient, 'start', element).then((resp) => {
      if (this._closed) {
        return
      }
      console.log('run with:', resp)
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been started'),
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
  onClickStop(element: Element) {
    this._disabled = true
    ServerAPI.v1.proxys.postOne(this.httpClient, 'stop', null).then(() => {
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
  onClickTestOne(element: Element) {
    this._disabled = true
    element.request = true
    element.error = undefined
    element.duration = undefined
    ServerAPI.v1.proxys.postOne<{
      duration: number,
      text: string,
      error?: string,
    }>(this.httpClient,
      'testOne',
      element.outbound,
    ).then((resp) => {
      if (this._closed) {
        return
      } else if (resp.error) {
        const e = resp.error
        console.log("test with:", resp.text)
        this.toasterService.pop('error',
          this.i18nService.get('error'),
          e,
        )
        element.error = e
        return
      }

      const data = resp.duration
      console.log("test with:", resp.text)
      if (isNumber(data)) {
        element.duration = data
      }
      this.toasterService.pop('success',
        this.i18nService.get('test speed success'),
        `delay ${element.duration} ms`,
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
      element.error = e
    }).finally(() => {
      element.request = undefined
      if (this._closed) {
        return
      }
      this._disabled = false
    })
  }
  onClickSetIPTables(element: Element) {
    this._disabled = true
    ServerAPI.v1.iptables.postOne(this.httpClient, 'init', element.outbound).then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('iptables has been init'),
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
  onClickRestoreIPTables(element: Element) {
    this._disabled = true
    ServerAPI.v1.iptables.postOne(this.httpClient, 'restore', element.outbound).then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('iptables has been restore'),
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
  onClickEdit(element: Element) {
    this.matDialog.open(EditComponent, {
      data: {
        panel: this.panel,
        element: element,
      },
      disableClose: true,
    })
  }
  onClickDelete(element: Element) {
    this.matDialog.open(ConfirmComponent, {
      data: {
        title: this.i18nService.get("delete proxy element title"),
        content: `${this.i18nService.get("delete proxy element")} : ${element.id} ${element.outbound.name}`,
      },
    }).afterClosed().toPromise().then((data) => {
      if (this._closed || !data) {
        return
      }
      this._delete(element)
    })
  }
  private _delete(element: Element) {
    this._disabled = true
    ServerAPI.v1.proxys.delete(this.httpClient, {
      params: {
        subscription: this.panel.id.toString(),
        id: element.id.toString(),
      },
    }).then(() => {
      if (this._closed) {
        return
      }
      const index = this.panel.source.indexOf(element)
      this.panel.source.splice(index, 1)
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been deleted'),
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
  onClickCopySource(items: Array<Element>) {
    try {
      const str = items.map((v) => v.toShare()).join("\n")
      this._btnClipboard.nativeElement.setAttribute("data-clipboard-text", str)
      this._btnClipboard.nativeElement.click()
    } catch (e: any) {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        e,
      )
    }
  }
  onClickShareSource(items: Array<Element>) {
    try {
      const str = items.map((v) => v.toShare()).join("\n")
      this.matDialog.open(QrcodeComponent, {
        data: str,
      })
    } catch (e: any) {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        e,
      )
    }
  }
  onClickCopy(element: Element) {
    try {
      const str = element.toShare()
      this._btnClipboard.nativeElement.setAttribute("data-clipboard-text", str)
      this._btnClipboard.nativeElement.click()
    } catch (e: any) {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        e,
      )
    }
  }
  onClickShare(element: Element) {
    try {
      const str = element.toShare()
      this.matDialog.open(QrcodeComponent, {
        data: str,
      })
    } catch (e: any) {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        e,
      )
    }
  }
}
