import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './app/home/home.component';
import { IptablesSaveComponent } from './app/iptables/iptables-save/iptables-save.component';
import { IptablesTemplateComponent } from './app/iptables/iptables-template/iptables-template.component';
const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
  },
  {
    path: 'content',
    loadChildren: () => import('./content/content.module').then(m => m.ContentModule),
  },
  {
    path: 'user',
    loadChildren: () => import('./user/user.module').then(m => m.UserModule),
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
    path: 'v2ray',
    loadChildren: () => import('./v2ray/v2ray.module').then(m => m.V2rayModule),
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
