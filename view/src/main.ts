import { MissingTranslationStrategy } from '@angular/core';
import { enableProdMode, TRANSLATIONS, TRANSLATIONS_FORMAT } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

import { AppModule } from './app/app.module';
import { environment } from './environments/environment';

import 'hammerjs';
declare const require: any
if (environment.production) {
  enableProdMode();
}

platformBrowserDynamic().bootstrapModule(AppModule)
  .catch(err => console.error(err));

let options = undefined;
if (!environment.production) {
  const translations = require(`raw-loader!./locale/zh-Hant.xlf`);
  options = {
    missingTranslation: MissingTranslationStrategy.Error,
    providers: [
      { provide: TRANSLATIONS, useValue: translations },
      { provide: TRANSLATIONS_FORMAT, useValue: 'xlf' }
    ]
  };
}

platformBrowserDynamic().bootstrapModule(AppModule, options)
  .catch(err => console.log(err));

