import { Pipe, PipeTransform } from '@angular/core';
import { AppService } from '../services/app.service';

@Pipe({
  name: 'commonText',
  pure: false,
})
export class CommonTextPipe implements PipeTransform {

  constructor(
    private appService: AppService,
  ) { }

  transform(value: any) {
    if (value === 'hours') {
      return this.appService.hoursName;
    } else if (value === 'coin') {
      return this.appService.coinName;
    } else if (value === 'coinFull') {
      return this.appService.fullCoinName;
    }

    return '';
  }
}
