import { Pipe, PipeTransform } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { BlockchainService } from '../services/blockchain.service';
import { AppService } from '../services/app.service';

@Pipe({
  name: 'amount',
  pure: false,
})
export class AmountPipe implements PipeTransform {

  constructor(
    private decimalPipe: DecimalPipe,
    private blockchainService: BlockchainService,
    private appService: AppService,
  ) { }

  transform(value: any, showingCoins = true, partToReturn = '') {
    let firstPart: string;
    let response = '';

    if (partToReturn !== 'last') {
      firstPart = this.decimalPipe.transform(value, showingCoins ? ('1.0-' + this.blockchainService.currentMaxDecimals) : '1.0-0');
      response = firstPart;
      if (partToReturn !== 'first') {
        response += ' ';
      }
    }
    if (partToReturn !== 'first') {
      response += showingCoins ? this.appService.coinName : this.appService.hoursName;
    }

    return response;
  }
}
