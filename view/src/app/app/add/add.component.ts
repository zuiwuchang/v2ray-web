import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Utils } from 'src/app/core/utils';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Panel, Outbound, Element } from '../view/source';
@Component({
  selector: 'app-add',
  templateUrl: './add.component.html',
  styleUrls: ['./add.component.scss']
})
export class AddComponent implements OnInit, OnDestroy {
  constructor(
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<AddComponent>,
    @Inject(MAT_DIALOG_DATA) public panel: Panel,
  ) { }

  private _closed = false
  private _disabled: boolean
  get disabled(): boolean {
    return this._disabled
  }
  outbound = new Outbound()
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this._closed = true
  }
  onClose() {
    this.matDialogRef.close()
  }
  onSave() {
    this._disabled = true
    this.outbound.format()
    this.httpClient.post<number>(ServerAPI.proxy.add, {
      subscription: this.panel.id,
      outbound: this.outbound,
    }).toPromise().then((id) => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element added successfully'),
      )
      const data = new Element()
      data.id = id
      data.subscription = this.panel.id
      data.outbound = this.outbound
      this.panel.source.push(data)
      this.panel.source.sort(Element.compare)
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