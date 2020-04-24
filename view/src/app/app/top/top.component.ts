import { Component, OnInit, OnDestroy } from '@angular/core';

import { StatusService, Status } from 'src/app/core/status/status.service';
import { Subscription } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { ServerAPI } from 'src/app/core/core/api';
import { Utils } from 'src/app/core/utils';
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
  ngOnInit(): void {
    this._subscription = this.statusService.observable.subscribe((status) => {
      if (this._closed) {
        return
      }
      this._status = status
    })
  }
  ngOnDestroy() {
    this._closed = true
    this._subscription.unsubscribe()
  }
  onClickStop() {
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
        Utils.resolveError(e),
      )
    }).finally(() => {
      this._disabled = false
    })
  }
}
