import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LicenseComponent } from './license/license.component';
import { AboutComponent } from './about/about.component';


const routes: Routes = [
  {
    path: 'license',
    component: LicenseComponent
  },
  {
    path: 'about',
    component: AboutComponent
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ContentRoutingModule { }
