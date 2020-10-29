import { Component, OnInit, VERSION } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { resolveError } from 'src/app/core/core/restful';
import { ToasterService } from 'angular2-toaster';
interface Version {
  platform: string
  tag: string
  commit: string
  date: string
  v2ray: string
}
@Component({
  selector: 'app-about',
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class AboutComponent implements OnInit {
  VERSION = VERSION
  content: string
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
  ) { }
  version: Version
  ngOnInit(): void {
    ServerAPI.v1.version.get<Version>(this.httpClient).then((data) => {
      this.version = data
    }, (e) => {
      console.warn(e)
    })
    this.httpClient.get(ServerAPI.static.licenses, {
      responseType: 'text',
    }).subscribe((text) => {
      this.content = text
    }, (e) => {

      this.toasterService.pop('error',
        undefined,
        resolveError(e),
      )
    })
  }

}
