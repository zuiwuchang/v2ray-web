import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ContentRoutingModule } from './content-routing.module';
import { AboutComponent } from './about/about.component';
import { LicenseComponent } from './license/license.component';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
@NgModule({
  declarations: [AboutComponent, LicenseComponent],
  imports: [
    CommonModule,
    ContentRoutingModule,
    MatCardModule, MatButtonModule, MatListModule,
    MatIconModule,
  ]
})
export class ContentModule { }
