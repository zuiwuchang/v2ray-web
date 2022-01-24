import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { V2rayRoutingModule } from './v2ray-routing.module';

import { SharedModule } from '../shared/shared.module';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatListModule } from '@angular/material/list';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';

import { SettingsComponent } from './settings/settings.component';
import { SubscriptionComponent } from './subscription/subscription.component';
import { SubscriptionAddComponent } from './subscription/subscription-add/subscription-add.component';
import { SubscriptionEditComponent } from './subscription/subscription-edit/subscription-edit.component';
import { PreviewComponent } from './dialog/preview/preview.component';


@NgModule({
    declarations: [SettingsComponent, SubscriptionComponent, SubscriptionAddComponent,
        SubscriptionEditComponent,
        PreviewComponent,
    ],
    imports: [
        CommonModule, FormsModule,
        V2rayRoutingModule,
        SharedModule,
        MatProgressBarModule, MatIconModule, MatButtonModule,
        MatListModule, MatDialogModule, MatFormFieldModule,
        MatInputModule, MatProgressSpinnerModule, MatTooltipModule,
        MatCardModule,
    ]
})
export class V2rayModule { }
