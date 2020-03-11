import { Pipe, PipeTransform } from '@angular/core';

import { AppService } from '../services/app.service';

/**
 * Returns the name of a commonly used element. The posible values are:
 * hours: returns the name of the coin hours.
 * coin: returns the short name of the coin, like 'SKY' for Skycoin.
 * coinFull: returns the full name of the coin, like 'Skycoin'.
 * The pipe expect the value to be exactly one of the previously listed strings.
 */
@Pipe({
  name: 'commonText',
  pure: false,
})
export class CommonTextPipe implements PipeTransform {

  constructor(
    private appService: AppService,
  ) { }

  transform(value: string) {
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
