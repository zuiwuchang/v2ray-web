import { Component, OnInit, Inject } from '@angular/core';
import {  MAT_DIALOG_DATA } from '@angular/material/dialog';
interface Data {
  text: string
  error: string
}
@Component({
  selector: 'app-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss']
})
export class PreviewComponent implements OnInit {

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: Data) { }

  ngOnInit(): void {
  }

}
