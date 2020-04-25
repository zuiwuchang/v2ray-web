import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { IptablesRoutingModule } from './iptables-routing.module';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';

import { TemplatesComponent } from './templates/templates.component';
import { ViewComponent } from './view/view.component';


@NgModule({
  declarations: [TemplatesComponent, ViewComponent],
  imports: [
    CommonModule, FormsModule,
    IptablesRoutingModule,
    MatButtonModule, MatProgressBarModule, MatCardModule,
    MatFormFieldModule, MatInputModule,
  ]
})
export class IptablesModule { }
