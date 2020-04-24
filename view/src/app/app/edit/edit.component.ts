import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Utils } from 'src/app/core/utils';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Panel, Outbound, Element } from '../view/source';

interface Data {
  panel: Panel
  element: Element
}
@Component({
  selector: 'app-edit',
  templateUrl: './edit.component.html',
  styleUrls: ['./edit.component.scss']
})
export class EditComponent implements OnInit, OnDestroy {

  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<EditComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Data, ) { }
  private _closed = false
  private _disabled: boolean
  get disabled(): boolean {
    return this._disabled
  }
  outbound = new Outbound()
  ngOnInit(): void {
    this.data.element.outbound.cloneTo(this.outbound)
  }
  ngOnDestroy() {
    this._closed = true
  }
  onClose() {
    this.matDialogRef.close()
  }
  get isNotChanged(): boolean {
    return this.data.element.outbound.equal(this.outbound)
  }
  onSave() {
    this._disabled = true
    this.outbound.format()
    this.httpClient.post<number>(ServerAPI.proxy.put, {
      id: this.data.element.id,
      subscription: this.data.panel.id,
      outbound: this.outbound,
    }).toPromise().then((id) => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element put successfully'),
      )
      this.outbound.cloneTo(this.data.element.outbound)
      this.matDialogRef.close()
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
