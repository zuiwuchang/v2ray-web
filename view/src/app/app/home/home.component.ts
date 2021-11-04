import { Component, OnInit, OnDestroy } from '@angular/core';
import { SessionService, Session } from 'src/app/core/session/session.service';
import { filter, first } from 'rxjs/operators';
import { Subscription } from 'rxjs';
var Value: any
@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit, OnDestroy {

  constructor(private sessionService: SessionService,
  ) { }
  private _closed = false
  private _session: Session = Value
  get session(): Session {
    return this._session
  }
  private _ready = false
  get ready(): boolean {
    return this._ready
  }
  private _subscription: Subscription = Value
  ngOnInit(): void {
    this.sessionService.ready.then(() => {
      this._ready = true
      if (this._closed) {
        return
      }
      this._subscription = this.sessionService.observable.pipe(
        filter(function (data) {
          if (data) {
            return true
          }
          return false
        }),
        first()
      ).subscribe((data) => {
        if (this._closed) {
          return
        }
        this._session = data
      })
    })
  }
  ngOnDestroy() {
    this._closed = true
    if (this._subscription) {
      this._subscription.unsubscribe()
    }
  }
}
