import { Component, Input, OnInit } from '@angular/core';
import { StrategyElement } from '../strategy'
@Component({
  selector: 'strategy-value',
  templateUrl: './value.component.html',
  styleUrls: ['./value.component.scss']
})
export class ValueComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }
  @Input('value')
  value = new StrategyElement()
  @Input('disabled')
  disabled = false
}
