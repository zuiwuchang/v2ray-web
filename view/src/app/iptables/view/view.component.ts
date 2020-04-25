import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { Utils } from 'src/app/core/utils';
import { isString } from 'util';
import { SessionService } from 'src/app/core/session/session.service';

@Component({
  selector: 'app-view',
  templateUrl: './view.component.html',
  styleUrls: ['./view.component.scss']
})
export class ViewComponent implements OnInit, OnDestroy {
  constructor(private httpClient: HttpClient,
    private sessionService: SessionService,
  ) { }
  private _ready = false
  get ready(): boolean {
    return this._ready
  }
  private _closed = false
  err: any
  text: string = ''
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
    this.httpClient.get<string>(ServerAPI.iptables.view).toPromise().then((text) => {
      if (this._closed) {
        return
      }
      this.text = text
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
}
