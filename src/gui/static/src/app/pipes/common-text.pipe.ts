import { Pipe, PipeTransform } from '@angular/core';
import { AppConfig } from '../app.config';

@Pipe({
  name: 'commonText',
})
export class CommonTextPipe implements PipeTransform {
  transform(value: any) {
    if (value === 'hours') {
      return AppConfig.hoursNamePlural;
    } else if (value === 'hour') {
      return AppConfig.hoursNameSingular;
    } else if (value === 'coin') {
      return AppConfig.coinName;
    } else if (value === 'coinFull') {
      return AppConfig.coinName;
    }

    return '';
  }
}
