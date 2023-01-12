import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { StrategyRoutingModule } from './strategy-routing.module';
import { ViewComponent } from './view/view.component';

import { MatListModule } from '@angular/material/list';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialogModule } from '@angular/material/dialog';
import { MatProgressBarModule } from '@angular/material/progress-bar';

@NgModule({
  declarations: [
    ViewComponent
  ],
  imports: [
    CommonModule,
    MatListModule,
    MatButtonModule,
    MatIconModule,
    MatTooltipModule,
    MatFormFieldModule,
    MatInputModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatProgressBarModule,
    StrategyRoutingModule
  ]
})
export class StrategyModule { }
