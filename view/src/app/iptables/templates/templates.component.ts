import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { ContextText } from '../../core/text';
import { SessionService } from 'src/app/core/session/session.service';
class Context {
  shell: string = ''
  view: string = ''
  clear: string = ''
  init: string = ''
  cloneFrom(other: Context) {
    this.shell = other.shell
    this.view = other.view
    this.clear = other.clear
    this.init = other.init
  }
}
@Component({
  selector: 'app-templates',
  templateUrl: './templates.component.html',
  styleUrls: ['./templates.component.scss']
})
export class TemplatesComponent implements OnInit, OnDestroy {
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
  context = new Context()
  contextText = ContextText
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
    ServerAPI.v1.iptables.get<Context>(this.httpClient).then((data) => {
      if (this._closed) {
        return
      }
      if (data) {
        this.context.cloneFrom(data)
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
    ServerAPI.v1.iptables.put(this.httpClient, this.context).then(() => {
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
  onClickResetDefault() {
    this._disabled = true
    ServerAPI.v1.iptables.getOne<Context>(this.httpClient, 'default').then((data) => {
      if (this._closed) {
        return
      }
      if (data) {
        this.context.cloneFrom(data)
      }
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
