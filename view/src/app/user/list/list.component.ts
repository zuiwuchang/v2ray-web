import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Utils } from 'src/app/core/utils';
import { User } from '../user';
import { MatDialog } from '@angular/material/dialog';
import { ConfirmComponent } from 'src/app/shared/dialog/confirm/confirm.component';
import { AddComponent } from '../add/add.component';
import { PasswordComponent } from '../password/password.component';
@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss']
})
export class ListComponent implements OnInit, OnDestroy {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialog: MatDialog,
  ) { }
  private _closed = false
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  private _source = new Array<User>()
  get source(): Array<User> {
    return this._source
  }
  ngOnInit(): void {
    this.httpClient.get<Array<User>>(ServerAPI.user.list).toPromise().then((data) => {
      if (this._closed) {
        return
      }
      if (data && data.length > 0) {
        this._source.push(...data)
        this._source.sort(User.compare)
      }
    }, (e) => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        Utils.resolveError(e),
      )
    })
  }
  ngOnDestroy(): void {
    this._closed = true
  }
  onClickEdit(node: User) {
    this.matDialog.open(PasswordComponent, {
      data: node.name,
    })
  }
  onClickDelete(node: User) {
    this.matDialog.open(ConfirmComponent, {
      data: {
        title: this.i18nService.get("delete user title"),
        content: `${this.i18nService.get("delete user")} : ${node.name}`,
      },
    }).afterClosed().toPromise().then((data) => {
      if (this._closed || !data) {
        return
      }
      this._delete(node)
    })
  }
  private _delete(node: User) {
    this._disabled = true
    this.httpClient.post(ServerAPI.user.remove, {
      name: node.name,
    }).toPromise().then(() => {
      const index = this._source.indexOf(node)
      this._source.splice(index, 1)
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('user has been deleted'),
      )
    }, (e) => {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        Utils.resolveError(e),
      )
    }).finally(() => {
      this._disabled = false
    })
  }
  onClickAdd() {
    this.matDialog.open(AddComponent, {
      disableClose: true,
    }).afterClosed().toPromise().then((data) => {
      if (this._closed || !data) {
        return
      }
      this._source.push(data)
      this._source.sort(User.compare)
    })
  }
}
