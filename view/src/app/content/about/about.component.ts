import { Component, OnInit, VERSION } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
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
  constructor(private httpClient: HttpClient,
  ) { }
  version: Version
  ngOnInit(): void {
    this.httpClient.get<Version>(ServerAPI.version).toPromise().then((data) => {
      this.version = data
    }, (e) => {
      console.warn(e)
    })
  }

}
