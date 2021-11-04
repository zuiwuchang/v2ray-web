import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Element } from '../subscription';
import { isString } from 'king-node/dist/core';
@Component({
  selector: 'app-subscription-edit',
  templateUrl: './subscription-edit.component.html',
  styleUrls: ['./subscription-edit.component.scss']
})
export class SubscriptionEditComponent implements OnInit, OnDestroy {

  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<SubscriptionEditComponent>,
    @Inject(MAT_DIALOG_DATA) private data: Element,
  ) { }
  private _closed = false
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  name: string = ''
  url: string = ''
  ngOnInit(): void {
    if (!isString(this.data.name)) {
      this.data.name = ''
    } else {
      this.data.name = this.data.name.trim()
    }
    if (!isString(this.data.url)) {
      this.data.url = ''
    } else {
      this.data.url = this.data.url.trim()
    }
    this.name = this.data.name
    this.url = this.data.url
  }
  ngOnDestroy() {
    this._closed = true
  }
  get isNotChanged(): boolean {
    return this.data.name.trim() == this.name.trim() && this.data.url.trim() == this.url.trim()
  }
  onSave() {
    this._disabled = true
    const name = this.name.trim()
    const url = this.url.trim()
    ServerAPI.v1.subscriptions.put(this.httpClient, {
      id: this.data.id,
      name: name,
      url: url,
    }).then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('subscription reset complete'),
      )
      this.matDialogRef.close({
        name: name,
        url: url,
      })
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
