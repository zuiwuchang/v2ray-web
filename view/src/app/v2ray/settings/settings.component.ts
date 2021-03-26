import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { isString } from 'king-node/dist/core';
import { ContextText } from '../../core/text';
import { SessionService } from 'src/app/core/session/session.service';
import { MatDialog } from '@angular/material/dialog';
import { PreviewComponent } from '../dialog/preview/preview.component';
interface Response {
  text: string
}
@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit, OnDestroy {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private sessionService: SessionService,
    private matDialog: MatDialog,
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
  url: string = 'vmess://eyJwcyI6InRlc3QiLCJhZGQiOiJmdWNrY2NwLmNvbSIsInBvcnQiOiI0NDMiLCJob3N0IjoiZnVja2NjcC5jb20iLCJ0bHMiOiJ0bHMiLCJuZXQiOiJ3cyIsInBhdGgiOiIvZnVja2NjcC9GaWdodE9yZGllIiwiaWQiOiIxZWYxOTdmNi03NzA4LTQ3NDItYjA5Zi1lNTBjNWVkMTVmNWUiLCJhaWQiOiIwIiwidHlwZSI6ImF1dG8iLCJ2IjoiMCJ9'
  contextText = ContextText
  private _text: string = ''
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
    ServerAPI.v1.v2ray.get<string>(this.httpClient).then((text) => {
      if (this._closed) {
        return
      }
      if (isString(text)) {
        this.text = text.trim()
      } else {
        this.text = ''
      }
      this._text = text
    }, (e) => {
      if (this._closed) {
        return
      }
      console.warn(e)
      this.err = e
    }).finally(() => {
      this._ready = true
    })
  }

  onClickSave() {
    this._disabled = true
    ServerAPI.v1.v2ray.put(this.httpClient, {
      text: this.text,
    }).then(() => {
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
        e,
      )
    }).finally(() => {
      this._disabled = false
    })
  }
  get isNotChange(): boolean {
    return this.text.trim() == this._text.trim()
  }
  onClickPreview() {
    if (this.url == "") {
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        this.i18nService.get('proxy url not support empty'),
      )
      return
    }
    this._disabled = true
    ServerAPI.v1.v2ray.postOne(this.httpClient, "preview", {
      text: this.text,
      url: this.url,
    }).then((obj) => {
      if (this._closed) {
        return
      }
      this.matDialog.open(PreviewComponent, {
        data: obj,
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
  onClickTest() {
    if (this.url == "") {
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        this.i18nService.get('proxy url not support empty'),
      )
      return
    }
    this._disabled = true
    ServerAPI.v1.v2ray.postOne(this.httpClient, "test", {
      text: this.text,
      url: this.url,
    }).then(() => {
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
        e,
      )
    }).finally(() => {
      this._disabled = false
    })
  }
  onClickResetDefault() {
    this._disabled = true
    ServerAPI.v1.v2ray.getOne<Response>(this.httpClient, 'default').then((data) => {
      if (this._closed) {
        return
      }
      this.text = data.text
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
}
