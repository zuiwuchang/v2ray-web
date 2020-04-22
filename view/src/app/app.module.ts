import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { ToasterModule, ToasterService } from 'angular2-toaster';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './app/home/home.component';

import { SharedModule } from './shared/shared.module';
import { IptablesSaveComponent } from './app/iptables/iptables-save/iptables-save.component';
import { IptablesTemplateComponent } from './app/iptables/iptables-template/iptables-template.component';
import { V2raySettingComponent } from './app/v2ray/v2ray-setting/v2ray-setting.component';
import { V2raySubscriptionComponent } from './app/v2ray/v2ray-subscription/v2ray-subscription.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    IptablesSaveComponent,
    IptablesTemplateComponent,
    V2raySettingComponent,
    V2raySubscriptionComponent,
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
