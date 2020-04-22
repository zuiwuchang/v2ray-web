import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { SettingsComponent } from './settings/settings.component';
import { SubscriptionComponent } from './subscription/subscription.component';
const routes: Routes = [
  {
    path: 'settings',
    component: SettingsComponent,
  },
  {
    path: 'subscription',
    component: SubscriptionComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class V2rayRoutingModule { }
