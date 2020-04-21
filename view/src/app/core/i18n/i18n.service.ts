import { Injectable } from '@angular/core';
import { isString } from 'util';

@Injectable({
  providedIn: 'root'
})
export class I18nService {
  private _keys = new Map<string, string>()
  constructor() { }
  init(doc: any) {
    if (!doc || !doc.childNodes || doc.childNodes.length == 0) {
      return
    }
    for (let i = 0; i < doc.childNodes.length; i++) {
      const item = doc.childNodes[i]
      if (!item.attributes) {
        continue
      }
      let key = item.attributes["data-key"]
      if (key == "" || key == undefined || key == null) {
        continue
      }
      key = key.value
      if (key == "" || key == undefined || key == null) {
        continue
      }
      this._keys.set(key, item.innerText)
    }
  }
  get(key: string) {
    const val = this._keys.get(key)
    if (isString(val)) {
      return val
    }
    return key
  }
}
