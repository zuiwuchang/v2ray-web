import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { TemplatesComponent } from './templates/templates.component';
import { ViewComponent } from './view/view.component';


const routes: Routes = [
  {
    path: 'view',
    component: ViewComponent,
  },
  {
    path: 'templates',
    component: TemplatesComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class IptablesRoutingModule { }
