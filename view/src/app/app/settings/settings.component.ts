import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { SessionService } from 'src/app/core/session/session.service';
import { ServerAPI } from 'src/app/core/core/api';
import { isString } from 'king-node/dist/core';
interface Result {
  url: string
  v2ray: boolean
  iptables: boolean
}
@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private sessionService: SessionService,
  ) { }
  private _ready = false
  get ready(): boolean {
    return this._ready
  }
  private _closed = false
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  err: any
  url: string
  v2ray: boolean
  iptables: boolean
  ngOnInit(): void {
    this.sessionService.ready.then(() => {
      if (this._closed) {
        return
      }
      this.load()
    })
  }
  ngOnDestroy() {
    this._closed = true
  }
  load() {
    this.err = null
    this._ready = false
    ServerAPI.v1.settings.get<Result>(this.httpClient).then((data) => {
      if (this._closed) {
        return
      }
      if (data) {
        if (isString(data.url)) {
          this.url = data.url
        }
        if (data.iptables) {
          this.iptables = true
        }
        if (data.v2ray) {
          this.v2ray = true
        }
      }
    }, (e) => {
      if (this._closed) {
        return
      }
      console.warn(e)
      this.err = e
    }).finally(() => {
      this._ready = true
    })
  }
  onClickSave() {
    this._disabled = true
    ServerAPI.v1.settings.put(this.httpClient, {
      url: this.url,
      v2ray: this.v2ray,
      iptables: this.iptables,
    }).then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('data saved'),
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
}
