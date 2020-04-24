import { Component, OnInit, Input, OnDestroy } from '@angular/core';
import { Panel, Element } from '../view/source';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Utils } from 'src/app/core/utils';
import { isArray } from 'util';
import { MatDialog } from '@angular/material/dialog';
import { AddComponent } from '../add/add.component';
import { EditComponent } from '../edit/edit.component';
import { ConfirmComponent } from 'src/app/shared/dialog/confirm/confirm.component';
@Component({
  selector: 'app-view-panel',
  templateUrl: './view-panel.component.html',
  styleUrls: ['./view-panel.component.scss']
})
export class ViewPanelComponent implements OnInit, OnDestroy {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialog: MatDialog,
  ) { }
  @Input('panel')
  panel: Panel
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  private _closed = false
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this._closed = true
  }
  onClickTest() {
    console.log('test')
  }
  onClickAdd() {
    this.matDialog.open(AddComponent, {
      data: this.panel,
      disableClose: true,
    })
  }
  onClickClear() {
    this.matDialog.open(ConfirmComponent, {
      data: {
        title: this.i18nService.get("clear proxy element title"),
        content: this.i18nService.get("clear proxy element"),
      },
    }).afterClosed().toPromise().then((data) => {
      if (this._closed || !data) {
        return
      }
      this._clear()
    })
  }
  private _clear() {
    this._disabled = true
    this.httpClient.post(ServerAPI.proxy.clear, {
      subscription: this.panel.id,
    }).toPromise().then(() => {
      this.panel.source = new Array<Element>()
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been cleared'),
      )
    }, (e) => {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        Utils.resolveError(e),
      )
    }).finally(() => {
      this._disabled = false
    })
  }
  onClickUpdate() {
    this._disabled = true
    this.httpClient.post<Array<any>>(ServerAPI.proxy.update, {
      id: this.panel.id,
    }).toPromise().then((data) => {
      if (this._closed) {
        return
      }
      const source = new Array<Element>()
      this.panel.source.slice(0, this.panel.source.length)
      if (isArray(data) && data.length > 0) {
        for (let i = 0; i < data.length; i++) {
          source.push(new Element(data[i]))
        }
        source.sort(Element.compare)
      }
      this.panel.source = source
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been updated'),
      )
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
  onClickStart(element: Element) {
    console.log("start", element)
  }
  onClickStop(element: Element) {
    console.log("stop", element)
  }
  onClickSetIPTables(element: Element) {
    console.log("set iptables", element)
  }
  onClickRestoreIPTables(element: Element) {
    console.log("restore iptables", element)
  }
  onClickEdit(element: Element) {
    this.matDialog.open(EditComponent, {
      data: {
        panel: this.panel,
        element: element,
      },
      disableClose: true,
    })
  }
  onClickDelete(element: Element) {
    this.matDialog.open(ConfirmComponent, {
      data: {
        title: this.i18nService.get("delete proxy element title"),
        content: `${this.i18nService.get("delete proxy element")} : ${element.id} ${element.outbound.name}`,
      },
    }).afterClosed().toPromise().then((data) => {
      if (this._closed || !data) {
        return
      }
      this._delete(element)
    })
  }
  private _delete(element: Element) {
    this._disabled = true
    this.httpClient.post(ServerAPI.proxy.remove, {
      subscription: this.panel.id,
      id: element.id,
    }).toPromise().then(() => {
      const index = this.panel.source.indexOf(element)
      this.panel.source.splice(index, 1)
      this.toasterService.pop('success',
        this.i18nService.get('success'),
        this.i18nService.get('proxy element has been deleted'),
      )
    }, (e) => {
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        Utils.resolveError(e),
      )
    }).finally(() => {
      this._disabled = false
    })
  }
}
