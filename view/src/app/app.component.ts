import { Component, ViewChild, ElementRef, AfterViewInit } from '@angular/core';
import { ToasterConfig } from 'angular2-toaster';
import { I18nService } from './core/i18n/i18n.service';
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements AfterViewInit {
  constructor(
    private i18nService: I18nService
  ) { }
  config: ToasterConfig =
    new ToasterConfig({
      positionClass: "toast-bottom-right"
    })
  @ViewChild("xi18n")
  private xi18nRef?: ElementRef
  ngAfterViewInit() {
    if (!this.xi18nRef) {
      return
    }
    this.i18nService.init(this.xi18nRef.nativeElement);
  }
}
