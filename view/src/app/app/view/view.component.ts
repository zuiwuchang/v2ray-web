import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { isArray } from 'king-node/dist/core';
import { Source, Panel, Element } from './source';
interface Response {
  subscription: [{
    id: number
    name: string
    url: string
  }
  ]
  element: Array<any>
}
@Component({
  selector: 'app-view',
  templateUrl: './view.component.html',
  styleUrls: ['./view.component.scss']
})
export class ViewComponent implements OnInit, OnDestroy {

  constructor(private httpClient: HttpClient,

  ) { }
  private _ready = false
  get ready(): boolean {
    return this._ready
  }
  private _closed = false
  err: any
  private _source = new Source()
  get source(): Source {
    return this._source
  }

  ngOnInit(): void {

    const panel = new Panel()
    panel.id = 0
    this.source.put(panel)
    this.load()
  }
  ngOnDestroy() {
    this._closed = true
  }
  load() {
    this.err = null
    this._ready = false
    ServerAPI.v1.proxys.get<Response>(this.httpClient).then((response) => {
      if (this._closed) {
        return
      }
      if (isArray(response.subscription)) {
        for (let i = 0; i < response.subscription.length; i++) {
          const element = response.subscription[i]
          const panel = new Panel()
          panel.id = element.id
          panel.name = element.name
          this._source.put(panel)
        }
      }
      if (isArray(response.element)) {
        for (let i = 0; i < response.element.length; i++) {
          const element = response.element[i]
          this._source.set(new Element(element))
        }
      }
      this._source.sort()
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
  onClick(evt: Event) {
    evt.stopPropagation()
  }
}
