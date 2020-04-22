import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { Utils } from 'src/app/core/utils';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { isString } from 'util';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit, OnDestroy {
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
  err: any
  text: string = ''
  private _text: string = ''
  ngOnInit(): void {
    this.load()
  }
  ngOnDestroy() {
    this._closed = true
  }
  async load() {
    this.err = null
    this._ready = false
    try {
      const text = await this.httpClient.get<string>(ServerAPI.v2ray.settings.get).toPromise()
      if (this._closed) {
        return
      }
      if (isString(text)) {
        this.text = text.trim()
      } else {
        this.text = ''
      }
      this._text = text
    } catch (e) {
      if (this._closed) {
        return
      }
      console.log(e)
      this.err = Utils.resolveError(e)
    } finally {
      this._ready = true
    }
  }

  onClickSave() {
    this._disabled = true
    this.httpClient.post(ServerAPI.v2ray.settings.put, {
      text: this.text,
    }).toPromise().then(() => {
      if (this._closed) {
        return
      }
      this._text = this.text
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
        Utils.resolveError(e),
      )
    }).finally(() => {
      this._disabled = false
    })
  }
  get isNotChange(): boolean {
    return this.text.trim() == this._text.trim()
  }
  onClickTest() {
    this._disabled = true
    this.httpClient.post(ServerAPI.v2ray.settings.test, {
      text: this.text,
    }).toPromise().then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('test success'),
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
  onClickResetDefault() {
    this.text = `{
    "log": {
        "loglevel": "warning"
    },
    "dns": {
        "servers": [
            // 使用 google 解析
            {
                "address": "8.8.8.8",
                "port": 53,
                "domains": [
                    "geosite:google",
                    "geosite:facebook",
                    "geosite:geolocation-!cn"
                ]
            },
            // 使用 趙國 解析服務
            {
                "address": "114.114.114.114",
                "port": 53,
                "domains": [
                    "geosite:cn",
                    "geosite:speedtest",
                    "domain:cn"
                ]
            },
            "8.8.8.8",
            "8.8.4.4",
            "localhost"
        ]
    },
    "inbounds": [
        // 本地 socks5 代理
        {
            "tag": "socks",
            "listen": "127.0.0.1",
            "protocol": "socks",
            "port": 1080,
            "settings": {
                "auth": "noauth"
            }
        },
        // 透明代理
        {
            "tag": "redir",
            "protocol": "dokodemo-door",
            "port": 10090,
            "settings": {
                "network": "tcp,udp",
                "followRedirect": true
            },
            "sniffing": {
                "enabled": true,
                "destOverride": [
                    "http",
                    "tls"
                ]
            }
        },
        // dns 代理 解決 域名污染
        {
            "tag": "dns",
            "protocol": "dokodemo-door",
            "port": 10054,
            "settings": {
                "address": "8.8.8.8",
                "port": 53,
                "network": "tcp,udp",
                "followRedirect": false
            }
        }
    ],
    "outbounds": [
        // 代理 訪問
        {
            "tag": "proxy",
            "protocol": "vmess",
            "settings": {
                "vnext": [
                    {{.Vnext}}
                ]
            },
            "streamSettings": {{.StreamSettings}},
            "mux": {
                "enabled": true
            }
        },
        // 直接 訪問
        {
            "tag": "freedom",
            "protocol": "freedom",
            "settings": {}
        },
        // 拒絕 訪問
        {
            "tag": "blackhole",
            "protocol": "blackhole",
            "settings": {}
        }
    ],
    "routing": {
        "domainStrategy": "IPIfNonMatch",
        "rules": [
            // 通過透明代理 進入 一律 代理訪問
            {
                "type": "field",
                "network": "tcp,udp",
                "inboundTag": [
                    "redir",
                    "dns"
                ],
                "outboundTag": "proxy"
            },
            // 代理訪問
            {
                "type": "field",
                "domain": [
                    "geosite:google",
                    "geosite:facebook",
                    "geosite:geolocation-!cn"
                ],
                "network": "tcp,udp",
                "outboundTag": "proxy"
            },
            // 直接訪問
            {
                "type": "field",
                "domain": [
                    "geosite:cn",
                    "geosite:speedtest",
                    "domain:cn",
                    "geoip:private"
                ],
                "ip": [
                    "geoip:cn"
                ],
                "network": "tcp,udp",
                "outboundTag": "freedom"
            }
        ]
    }
}`
  }
}
