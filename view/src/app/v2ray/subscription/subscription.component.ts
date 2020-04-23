import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { Utils } from 'src/app/core/utils';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { SessionService } from 'src/app/core/session/session.service';
import { isArray } from 'util';
import { Element } from './subscription';
import { MatDialog } from '@angular/material/dialog';
import { SubscriptionAddComponent } from './subscription-add/subscription-add.component';
import { ConfirmComponent } from 'src/app/shared/dialog/confirm/confirm.component';
import { SubscriptionEditComponent } from './subscription-edit/subscription-edit.component';

@Component({
  selector: 'app-subscription',
  templateUrl: './subscription.component.html',
  styleUrls: ['./subscription.component.scss']
})
export class SubscriptionComponent implements OnInit, OnDestroy {

  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialog: MatDialog,
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
  private _source = new Array<Element>()
  get source(): Array<Element> {
    return this._source
  }
  err: any
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
    this.httpClient.get<Array<Element>>(ServerAPI.v2ray.subscription.list).toPromise().then((source) => {
      if (this._closed) {
        return
      }
      if (isArray(source) && source.length > 0) {
        this._source.push(...source)
        this._source.sort(Element.compare)
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
  onClickAdd() {
    this.matDialog.open(SubscriptionAddComponent, {
      disableClose: true,
    }).afterClosed().toPromise().then((data) => {
      if (this._closed || !data) {
        return
      }
      this._source.push(data)
      this._source.sort(Element.compare)
    })
  }
  onClickEdit(element: Element) {
    this.matDialog.open(SubscriptionEditComponent, {
      data: {
        id: element.id,
        name: element.name,
        url: element.url,
      },
      disableClose: true,
    }).afterClosed().toPromise().then((data) => {
      if (this._closed || !data) {
        return
      }
      element.name = data.name
      element.url = data.url
    })
  }
  onClickDelete(element: Element) {
    this.matDialog.open(ConfirmComponent, {
      data: {
        title: this.i18nService.get("delete subscription title"),
        content: `${this.i18nService.get("delete subscription")} : ${element.name}`,
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
    this.httpClient.post(ServerAPI.v2ray.subscription.remove, {
      id: element.id,
    }).toPromise().then(() => {
      const index = this._source.indexOf(element)
      this._source.splice(index, 1)
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('subscription has been deleted'),
      )
    }, (e) => {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        Utils.resolveError(e),
      )
    }).finally(() => {
      this._disabled = false
    })
  }
}
