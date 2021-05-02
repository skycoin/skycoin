import { Pipe, PipeTransform } from '@angular/core';
import { BigNumber } from 'bignumber.js';

import { AppService } from '../services/app.service';

/**
 * Converts a number into a coin or hour amount. The resulting string is formatted with the
 * correct max number of decimals and the small name of the coin or the full name of the
 * coin hours at the end. It expects a number, numeric string or BigNumber instance as argument.
 * Also, the pipe can receive 2 optional arguments: the first one is a boolean value indicating
 * if the provided number must be converted to a coin amount (true) or an hour amount (false);
 * while the second one is a string which can be 'first', for the pipe to return only the first part
 * of the amount (the formatted number without the coin name), or 'last', for the pipe to return
 * only the last part of the amount (the coin or hours name).
 */
@Pipe({
  name: 'amount',
  pure: false,
})
export class AmountPipe implements PipeTransform {

  constructor(
    private appService: AppService,
  ) { }

  transform(value: any, showingCoins = true, partToReturn = '') {
    const convertedVal = new BigNumber(value).decimalPlaces(showingCoins ? this.appService.currentMaxDecimals : 0);

    let response = '';

    // Add the numeric part.
    if (partToReturn !== 'last') {
      if (convertedVal.isNaN()) {
        response = 'NaN';
      } else {
        response = convertedVal.toFormat();
      }

      if (partToReturn !== 'first') {
        response += ' ';
      }
    }

    // Add the name.
    if (partToReturn !== 'first') {
      response += showingCoins ? this.appService.coinName : (convertedVal.absoluteValue().isEqualTo(1) ? this.appService.hoursNameSingular : this.appService.hoursName);
    }

    return response;
  }
}
