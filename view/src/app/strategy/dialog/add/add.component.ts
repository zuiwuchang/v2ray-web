import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { Strategy, StrategyElement, StrategyValue } from '../../strategy';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { ServerAPI } from 'src/app/core/core/api';

@Component({
  selector: 'app-add',
  templateUrl: './add.component.html',
  styleUrls: ['./add.component.scss']
})
export class AddComponent implements OnInit {

  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<AddComponent>,
  ) { }
  private _closed = false
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  group = "Host"
  strategy = new StrategyValue({
    name: 'New',
    host: "dns.google, 8.8.8.8, 8.8.4.4\n",
    proxy: new StrategyElement({
      ip: "# 8.8.8.8, 8.8.4.4\n",
      domain: "# geosite:geolocation-!cn\n"
    }),
    direct: new StrategyElement({
      ip: "# geoip:private, geoip:cn\n",
      domain: "# geosite:cn\n"
    }),
    block: new StrategyElement({
      domain: "# geosite:category-ads-all\n# geosite:category-all\n"
    }),
  })
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this._closed = true
  }
  onSave() {
    this._disabled = true
    const strategy = this.strategy
    ServerAPI.v1.strategys.post(this.httpClient, {
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
        this.i18nService.get('strategy added successfully'),
      )
      this.matDialogRef.close(strategy)
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
