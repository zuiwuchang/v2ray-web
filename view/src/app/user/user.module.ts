import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../shared/shared.module';

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

import { UserRoutingModule } from './user-routing.module';
import { ListComponent } from './list/list.component';
import { AddComponent } from './add/add.component';
import { PasswordComponent } from './password/password.component';


@NgModule({
    declarations: [ListComponent, AddComponent, PasswordComponent],
    imports: [
        CommonModule, RouterModule, FormsModule,
        SharedModule,
        MatListModule, MatButtonModule, MatIconModule,
        MatTooltipModule, MatFormFieldModule, MatInputModule,
        MatCardModule, MatProgressSpinnerModule, MatDialogModule,
        MatProgressBarModule,
        UserRoutingModule
    ]
})
export class UserModule { }
