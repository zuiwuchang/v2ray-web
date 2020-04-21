import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './app/home/home.component';
import { AboutComponent } from './app/about/about.component';
import { LicenseComponent } from './app/license/license.component';
import { UserComponent } from './app/user/user.component';
import { IptablesSaveComponent } from './app/iptables/iptables-save/iptables-save.component';
import { IptablesTemplateComponent } from './app/iptables/iptables-template/iptables-template.component';
import { V2raySettingComponent } from './app/v2ray/v2ray-setting/v2ray-setting.component';
import { V2raySubscriptionComponent } from './app/v2ray/v2ray-subscription/v2ray-subscription.component';

const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
  },
  {
    path: 'about',
    component: AboutComponent,
  },
  {
    path: 'license',
    component: LicenseComponent,
  },
  {
    path: 'user',
    component: UserComponent,
  },
  {
    path: 'iptables/template',
    component: IptablesTemplateComponent,
  },
  {
    path: 'iptables/save',
    component: IptablesSaveComponent,
  },
  {
    path: 'v2ray/settings',
    component: V2raySettingComponent,
  },
  {
    path: 'v2ray/subscription',
    component: V2raySubscriptionComponent,
  },

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
