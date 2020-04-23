import { Directive, Input } from '@angular/core';
import { NG_VALIDATORS, AbstractControl } from '@angular/forms'
import { isString } from 'util';
@Directive({
  selector: '[sharedUrlValidator]',
  providers: [{ provide: NG_VALIDATORS, useExisting: UrlValidatorDirective, multi: true }]
})
export class UrlValidatorDirective {

  validate(control: AbstractControl): { [key: string]: any } {
    let str: string = control.value
    if (isString(str)) {
      str = str.trim()
      if (str.startsWith("http://") || str.startsWith("https://")) {
        return null
      }
    }

    return { 'url-validator': { value: control.value } }
  }
}
