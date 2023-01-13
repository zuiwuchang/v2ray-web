import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

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
import { MatSelectModule } from '@angular/material/select';
import { AddComponent } from './dialog/add/add.component';
import { ValueComponent } from './value/value.component';
import { EditComponent } from './dialog/edit/edit.component';

@NgModule({
  declarations: [
    ViewComponent,
    AddComponent,
    ValueComponent,
    EditComponent
  ],
  imports: [
    CommonModule, FormsModule,
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
    MatSelectModule,
    StrategyRoutingModule
  ]
})
export class StrategyModule { }
