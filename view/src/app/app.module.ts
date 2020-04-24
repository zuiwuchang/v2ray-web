import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { ToasterModule, ToasterService } from 'angular2-toaster';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule } from '@angular/material/list';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialogModule } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatMenuModule } from '@angular/material/menu';
import { MatDividerModule } from '@angular/material/divider';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './app/home/home.component';

import { SharedModule } from './shared/shared.module';
import { IptablesSaveComponent } from './app/iptables/iptables-save/iptables-save.component';
import { IptablesTemplateComponent } from './app/iptables/iptables-template/iptables-template.component';
import { ViewComponent } from './app/view/view.component';
import { ViewPanelComponent } from './app/view-panel/view-panel.component';
import { AddComponent } from './app/add/add.component';
import { EditComponent } from './app/edit/edit.component';
import { TopComponent } from './app/top/top.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    IptablesSaveComponent,
    IptablesTemplateComponent,
    ViewComponent,
    ViewPanelComponent,
    AddComponent,
    EditComponent,
    TopComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule, HttpClientModule, FormsModule,

    MatProgressBarModule, MatIconModule, MatButtonModule,
    MatListModule, MatExpansionModule, MatTooltipModule,
    MatDialogModule, MatProgressSpinnerModule, MatFormFieldModule,
    MatInputModule, MatAutocompleteModule, MatMenuModule,
    MatDividerModule,

    SharedModule,
    AppRoutingModule,
    ToasterModule.forRoot(),
  ],
  providers: [ToasterService],
  entryComponents: [AddComponent, EditComponent],
  bootstrap: [AppComponent]
})
export class AppModule { }
