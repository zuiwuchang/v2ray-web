import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { isString } from 'king-node/dist/core';
import { ContextText, V2rayTemplate } from '../../core/text';
import { SessionService } from 'src/app/core/session/session.service';
@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit, OnDestroy {
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
  text: string = ''
  contextText = ContextText
  private _text: string = ''
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
    this.httpClient.get<string>(ServerAPI.v2ray.settings.get).toPromise().then((text) => {
      if (this._closed) {
        return
      }
      if (isString(text)) {
        this.text = text.trim()
      } else {
        this.text = ''
      }
      this._text = text
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
    this.httpClient.post(ServerAPI.v2ray.settings.put, {
      text: this.text,
    }).toPromise().then(() => {
      if (this._closed) {
        return
      }
      this._text = this.text
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
  get isNotChange(): boolean {
    return this.text.trim() == this._text.trim()
  }
  onClickTest() {
    this._disabled = true
    this.httpClient.post(ServerAPI.v2ray.settings.test, {
      text: this.text,
    }).toPromise().then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('test success'),
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
  onClickResetDefault() {
    this.text = V2rayTemplate
  }
}
