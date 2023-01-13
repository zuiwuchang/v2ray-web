import { HttpClient } from '@angular/common/http';
import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { StrategyValue } from '../../strategy'

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
    @Inject(MAT_DIALOG_DATA) private data: StrategyValue,
  ) {
  }
  group = "Host"
  strategy = new StrategyValue()
  private _closed = false
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  ngOnInit(): void {
    this.data.cloneTo(this.strategy)
  }
  ngOnDestroy() {
    this._closed = true
  }
  onSave() {
    this._disabled = true
    const strategy = this.strategy
    ServerAPI.v1.strategys.put(this.httpClient, {
      name: strategy.name,
      value: strategy.value,
      host: strategy.host,
      proxyIP: strategy.proxy.ip,
      proxyDomain: strategy.proxy.domain,
      directIP: strategy.direct.ip,
      directDomain: strategy.direct.domain,
      blockIP: strategy.block.ip,
      blockDomain: strategy.block.domain,
    }).then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('strategy has been saved'),
      )
      strategy.cloneTo(this.data)
      this.matDialogRef.close()
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
