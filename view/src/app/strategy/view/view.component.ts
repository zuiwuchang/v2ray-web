import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { SessionService } from 'src/app/core/session/session.service';
import { sortNameValue } from 'src/app/core/utils';
import { Strategy } from '../strategy';

@Component({
  selector: 'app-view',
  templateUrl: './view.component.html',
  styleUrls: ['./view.component.scss']
})
export class ViewComponent implements OnInit {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialog: MatDialog,
    private sessionService: SessionService,
  ) { }
  private _ready = false
  get ready(): boolean {
    return this._ready
  }
  err: any
  private _closed = false
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  private _source = new Array<Strategy>()
  get source(): Array<Strategy> {
    return this._source
  }
  ngOnInit(): void {
    this.sessionService.ready.then(() => {
      if (this._closed) {
        return
      }
      this.load()
    })
  }
  ngOnDestroy(): void {
    this._closed = true
  }
  load() {
    this.err = null
    this._ready = false
    ServerAPI.v1.strategys.get<Array<Strategy>>(this.httpClient).then((data) => {
      if (this._closed) {
        return
      }
      if (data && data.length > 0) {
        this._source.push(...data)
        this._source.sort(sortNameValue)
      }
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
  onClickEdit(node: Strategy) {
    //   this.matDialog.open(PasswordComponent, {
    //     data: node.name,
    //     disableClose: true,
    //   })
  }
  onClickDelete(node: Strategy) {
    //   this.matDialog.open(ConfirmComponent, {
    //     data: {
    //       title: this.i18nService.get("delete user title"),
    //       content: `${this.i18nService.get("delete user")} : ${node.name}`,
    //     },
    //   }).afterClosed().toPromise().then((data) => {
    //     if (this._closed || !data) {
    //       return
    //     }
    //     this._delete(node)
    //   })
  }
  // private _delete(node: Strategy) {
  //   this._disabled = true
  //   ServerAPI.v1.users.delete(this.httpClient, {
  //     params: {
  //       name: node.name,
  //     },
  //   }).then(() => {
  //     const index = this._source.indexOf(node)
  //     this._source.splice(index, 1)
  //     this.toasterService.pop('success',
  //       this.i18nService.get('success'),
  //       this.i18nService.get('user has been deleted'),
  //     )
  //   }, (e) => {
  //     console.warn(e)
  //     this.toasterService.pop('error',
  //       this.i18nService.get('error'),
  //       e,
  //     )
  //   }).finally(() => {
  //     this._disabled = false
  //   })
  // }
  onClickAdd() {
    //   this.matDialog.open(AddComponent, {
    //     disableClose: true,
    //   }).afterClosed().toPromise().then((data) => {
    //     if (this._closed || !data) {
    //       return
    //     }
    //     this._source.push(data)
    //     this._source.sort(User.compare)
    //   })
  }
}
