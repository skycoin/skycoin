import { Pipe, PipeTransform } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { ExplorerService } from '../services/explorer/explorer.service';

/**
 * Formats a coin or coin hour amount. Normally, when using this pipe with a number, it limits
 * the number of decimal places and shows the short name of coin (or the long name of the coin
 * hours) after the amount, using the adecuate singular or plural form. However, it can be
 * configured in different ways.
 *
 *
 * Received value:
 *
 * The value must be a number. valid data types are "string" and "number".
 *
 *
 * Optional params:
 *
 * The first param is a boolean value that indicates if the amount to show is coins (true) or
 * coin hours (false). The default value is true.
 *
 * The second param is a string that can be used to limit what part of the amount will be shown.
 * If set to 'first', only the numeric part will be shown. If set to 'last', only the short name
 * of the coin or the long name of the coin hours (depending on the value of the first parameter
 * and using the numeric value to decide if the singular or plural form must be used) will be
 * shown. The default value is an empty string, so both parts are shown.
 */
@Pipe({
  name: 'amount',
  pure: false,
})
export class AmountPipe implements PipeTransform {

  constructor(
    private decimalPipe: DecimalPipe,
    private explorerService: ExplorerService,
  ) { }

  // Main function for processing the data.
  transform(value: any, showingCoins = true, partToReturn = '') {
    let firstPart: string;
    let response = '';

    if (partToReturn !== 'last') {
      // Use the standard decimalPipe for limiting the decimal places. The max number of decimal
      // places for the coins is obtained from the node. Coin hours never have decimal places.
      firstPart = this.decimalPipe.transform(value, showingCoins ? ('1.0-' + this.explorerService.maxDecimals) : '1.0-0');
      response = firstPart;
      if (partToReturn !== 'first') {
        response += ' ';
      }
    }
    if (partToReturn !== 'first') {
      response += showingCoins ? this.explorerService.coinName : (Number(value) === 1 || Number(value) === -1 ? this.explorerService.hoursNameSingular : this.explorerService.hoursName);
    }

    return response;
  }
}
