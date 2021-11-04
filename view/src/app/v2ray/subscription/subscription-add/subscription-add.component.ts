import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { MatDialogRef } from '@angular/material/dialog';
import { Element } from '../subscription';

@Component({
  selector: 'app-subscription-add',
  templateUrl: './subscription-add.component.html',
  styleUrls: ['./subscription-add.component.scss']
})
export class SubscriptionAddComponent implements OnInit, OnDestroy {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<SubscriptionAddComponent>,
  ) { }
  private _closed = false
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  name: string = ''
  url: string = ''
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this._closed = true
  }
  onSave() {
    this._disabled = true
    const name = this.name.trim()
    const url = this.url.trim()
    ServerAPI.v1.subscriptions.post<number>(this.httpClient, {
      name: name,
      url: url,
    }).then((id) => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('subscription added successfully'),
      )
      const data = new Element()
      data.id = id
      data.name = name
      data.url = url
      this.matDialogRef.close(data)
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
  onClose() {
    this.matDialogRef.close()
  }
}
