import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
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
import { MatCardModule } from '@angular/material/card';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './app/home/home.component';

import { SharedModule } from './shared/shared.module';
import { ViewComponent } from './app/view/view.component';
import { ViewPanelComponent } from './app/view-panel/view-panel.component';
import { AddComponent } from './app/add/add.component';
import { EditComponent } from './app/edit/edit.component';
import { TopComponent } from './app/top/top.component';
import { SettingsComponent } from './app/settings/settings.component';
import { QrcodeComponent } from './app/dialog/qrcode/qrcode.component';
import { ServiceWorkerModule } from '@angular/service-worker';
import { environment } from '../environments/environment';
import { HeaderInterceptor } from './app/service/header.interceptor';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    ViewComponent,
    ViewPanelComponent,
    AddComponent,
    EditComponent,
    TopComponent,
    SettingsComponent,
    QrcodeComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule, HttpClientModule, FormsModule,

    MatProgressBarModule, MatIconModule, MatButtonModule,
    MatListModule, MatExpansionModule, MatTooltipModule,
    MatDialogModule, MatProgressSpinnerModule, MatFormFieldModule,
    MatInputModule, MatAutocompleteModule, MatMenuModule,
    MatDividerModule, MatCardModule, MatCheckboxModule,
    MatSlideToggleModule,

    SharedModule,
    AppRoutingModule,
    ToasterModule.forRoot(),
    ServiceWorkerModule.register('ngsw-worker.js', { enabled: environment.production }),
  ],
  providers: [{
    provide: HTTP_INTERCEPTORS,
    useClass: HeaderInterceptor,
    multi: true,
  },
    ToasterService],
  entryComponents: [AddComponent, EditComponent, QrcodeComponent],
  bootstrap: [AppComponent]
})
export class AppModule { }
