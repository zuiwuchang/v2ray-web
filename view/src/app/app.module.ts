import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { ToasterModule, ToasterService } from 'angular2-toaster';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule } from '@angular/material/list';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTooltipModule } from '@angular/material/tooltip';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './app/home/home.component';

import { SharedModule } from './shared/shared.module';
import { IptablesSaveComponent } from './app/iptables/iptables-save/iptables-save.component';
import { IptablesTemplateComponent } from './app/iptables/iptables-template/iptables-template.component';
import { ViewComponent } from './app/view/view.component';
import { ViewPanelComponent } from './app/view-panel/view-panel.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    IptablesSaveComponent,
    IptablesTemplateComponent,
    ViewComponent,
    ViewPanelComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule, HttpClientModule,

    MatProgressBarModule, MatIconModule, MatButtonModule,
    MatListModule, MatExpansionModule, MatTooltipModule,

    SharedModule,
    AppRoutingModule,
    ToasterModule.forRoot(),
  ],
  providers: [ToasterService],
  bootstrap: [AppComponent]
})
export class AppModule { }
