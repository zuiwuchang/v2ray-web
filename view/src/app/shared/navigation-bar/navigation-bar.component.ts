import { Component, OnInit, OnDestroy } from '@angular/core';
import { Session, SessionService } from 'src/app/core/session/session.service';
import { Subscription } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { LoginComponent } from '../login/login.component';
@Component({
  selector: 'shared-navigation-bar',
  templateUrl: './navigation-bar.component.html',
  styleUrls: ['./navigation-bar.component.scss'],
})
export class NavigationBarComponent implements OnInit, OnDestroy {
  constructor(
    private sessionService: SessionService,
    private matDialog: MatDialog
  ) {}
  private _ready = false;
  get ready(): boolean {
    return this._ready;
  }
  private _session: Session;
  get session(): Session {
    return this._session;
  }
  private _subscription: Subscription;
  ngOnInit(): void {
    this.sessionService.ready.then((data) => {
      this._ready = data;
    });
    this._subscription = this.sessionService.observable.subscribe((data) => {
      this._session = data;
    });
  }
  ngOnDestroy(): void {
    this._subscription.unsubscribe();
  }

  onClickLogin() {
    this.matDialog.open(LoginComponent);
  }
  onClickLogout() {
    this.sessionService.logout();
  }
}
