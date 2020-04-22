import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { V2rayRoutingModule } from './v2ray-routing.module';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';

import { SettingsComponent } from './settings/settings.component';
import { SubscriptionComponent } from './subscription/subscription.component';


@NgModule({
  declarations: [SettingsComponent, SubscriptionComponent],
  imports: [
    CommonModule, FormsModule,
    V2rayRoutingModule,
    MatProgressBarModule, MatIconModule, MatButtonModule,
    MatCardModule,
  ]
})
export class V2rayModule { }
