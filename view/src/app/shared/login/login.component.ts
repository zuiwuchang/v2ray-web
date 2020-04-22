import { Component, OnInit, OnDestroy } from '@angular/core';
import { Session, SessionService } from 'src/app/core/session/session.service';
import { Subscription } from 'rxjs';
import { MatDialogRef } from '@angular/material/dialog';
import { Utils } from 'src/app/core/utils';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { sha512 } from 'js-sha512';
@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit, OnDestroy {
  constructor(private sessionService: SessionService,
    private matDialogRef: MatDialogRef<LoginComponent>,
    private toasterService: ToasterService,
    private i18nService: I18nService,
  ) { }
  private _ready = false
  get ready(): boolean {
    return this._ready
  }
  private _session: Session
  get session(): Session {
    return this._session
  }
  private _subscription: Subscription
  private _disabled: boolean
  get disabled(): boolean {
    return this._disabled
  }
  name: string
  password: string
  remember = true
  visibility = false
  ngOnInit(): void {
    this.sessionService.ready.then((data) => {
      this._ready = data
    })
    this._subscription = this.sessionService.observable.subscribe((data) => {
      this._session = data
    })
  }
  ngOnDestroy(): void {
    this._subscription.unsubscribe()
  }
  onClose() {
    this.matDialogRef.close()
  }
  async onSubmit() {
    try {
      const password = sha512(this.password).toString()
      this._disabled = true
      await this.sessionService.login(this.name, password, this.remember)
    } catch (e) {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        Utils.resolveError(e),
      )
    } finally {
      this._disabled = false
    }
  }
}
