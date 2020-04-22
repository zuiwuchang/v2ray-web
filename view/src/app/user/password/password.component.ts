import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Utils } from 'src/app/core/utils';
import { sha512 } from 'js-sha512';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { User } from '../user';
@Component({
  selector: 'app-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.scss']
})
export class PasswordComponent implements OnInit, OnDestroy {

  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<PasswordComponent>,
    @Inject(MAT_DIALOG_DATA) public name: string,
  ) {
  }
  private _closed = false
  private _disabled: boolean
  get disabled(): boolean {
    return this._disabled
  }
  password: string
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this._closed = true
  }
  onSave() {
    this._disabled = true
    const password = sha512(this.password).toString()
    this.httpClient.post(ServerAPI.user.password, {
      name: this.name,
      password: password,
    }).toPromise().then(() => {
      if (this._closed) {
        return
      }
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('password reset complete'),
      )
      const data = new User()
      data.name = this.name
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
  onClose() {
    this.matDialogRef.close()
  }
}
