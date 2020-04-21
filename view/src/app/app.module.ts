import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { ToasterModule, ToasterService } from 'angular2-toaster';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './app/home/home.component';
import { AboutComponent } from './app/about/about.component';
import { LicenseComponent } from './app/license/license.component';

import { SharedModule } from './shared/shared.module';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    AboutComponent,
    LicenseComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule, HttpClientModule,

    SharedModule,
    AppRoutingModule,
    ToasterModule.forRoot(),
  ],
  providers: [ToasterService],
  bootstrap: [AppComponent]
})
export class AppModule { }
