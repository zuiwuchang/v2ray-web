import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { ServerAPI } from 'src/app/core/core/api';
import { resolveError } from 'src/app/core/core/restful';

@Component({
  selector: 'app-license',
  templateUrl: './license.component.html',
  styleUrls: ['./license.component.scss']
})
export class LicenseComponent implements OnInit {

  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
  ) { }

  ngOnInit(): void {
    this.httpClient.get(
      ServerAPI.static.license,
      {
        responseType: 'text',
      },
    ).toPromise().then((data) => {
      this.content = data
    }, (e) => {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        resolveError(e),
      )
    })
  }
  content: any
}
