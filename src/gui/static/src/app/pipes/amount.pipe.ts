import { Pipe, PipeTransform } from '@angular/core';
import { AppConfig } from '../app.config';
import { DecimalPipe } from '@angular/common';
import { BlockchainService } from '../services/blockchain.service';

@Pipe({
  name: 'amount',
})
export class AmountPipe implements PipeTransform {

  constructor(
    private decimalPipe: DecimalPipe,
    private blockchainService: BlockchainService,
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
      response += showingCoins ? AppConfig.coinName : (firstPart !== '1' ? AppConfig.hoursNamePlural : AppConfig.hoursNameSingular);
    }

    return response;
  }
}
