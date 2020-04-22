import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { Utils } from 'src/app/core/utils';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { isString, isArray } from 'util';
import { Node } from './subscription';
@Component({
  selector: 'app-subscription',
  templateUrl: './subscription.component.html',
  styleUrls: ['./subscription.component.scss']
})
export class SubscriptionComponent implements OnInit, OnDestroy {

  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
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
  private _source = new Array<Node>()
  get source(): Array<Node> {
    return this._source
  }
  err: any
  ngOnInit(): void {
    this.load()
  }
  ngOnDestroy() {
    this._closed = true
  }
  load() {
    this.err = null
    this._ready = false
    this.httpClient.get<Array<Node>>(ServerAPI.v2ray.subscription.list).toPromise().then((source) => {
      if (this._closed) {
        return
      }
      if (isArray(source) && source.length > 0) {
        this._source.push(...source)
      }
    }, (e) => {
      if (this._closed) {
        return
      }
      console.warn(e)
      this.err = Utils.resolveError(e)
    }).finally(() => {
      this._ready = true
    })
  }
}
