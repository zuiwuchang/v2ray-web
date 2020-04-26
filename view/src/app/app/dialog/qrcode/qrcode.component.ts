import { Component, OnInit, Inject, ViewChild, ElementRef, AfterViewInit } from '@angular/core';
import * as QRCode from 'qrcode';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
@Component({
  selector: 'app-qrcode',
  templateUrl: './qrcode.component.html',
  styleUrls: ['./qrcode.component.scss']
})
export class QrcodeComponent implements OnInit, AfterViewInit {

  constructor(@Inject(MAT_DIALOG_DATA) private data: string) { }

  ngOnInit(): void {
  }
  @ViewChild("canvas")
  private _canvas: ElementRef
  ngAfterViewInit() {
    QRCode.toCanvas(this._canvas.nativeElement,
      this.data,
      (e) => {
        if (e) {
          console.warn(e)
        }
      })
  }
}
