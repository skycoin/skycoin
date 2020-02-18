import { filter } from 'rxjs/operators';
import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { SubscriptionLike } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { PriceService } from '../../../services/price.service';
import { BlockchainService } from '../../../services/blockchain.service';
import { AppService } from '../../../services/app.service';
import { NetworkService } from '../../../services/network.service';
import { AppConfig } from '../../../app.config';
import { BalanceAndOutputsService } from '../../../services/wallet-operations/balance-and-outputs.service';
import { AddressWithBalance } from '../../../services/wallet-operations/wallet-objects';

/**
 * Header shown at the top of most pages.
 */
@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Input() headline: string;

  // Data about the synchronization status of the node.
  synchronizationInfoObtained = false;
  synchronizationPercentage: number;
  // Use synchronizationInfoObtained to know if the value has been already updated.
  synchronized = false;
  currentBlock: number;
  highestBlock: number;

  // Data about the balance.
  coins: string;
  hours: string;

  showPrice = !!AppConfig.priceApiId;
  price: number;
  // If the node has pending transactions potentially affecting the user balance.
  hasPendingTxs: boolean;
  // If the app already got the balance from the node.
  balanceObtained = false;
  walletDownloadUrl = AppConfig.walletDownloadUrl;

  private subscriptionsGroup: SubscriptionLike[] = [];

  constructor(
    public appService: AppService,
    public networkService: NetworkService,
    private blockchainService: BlockchainService,
    private priceService: PriceService,
    private balanceAndOutputsService: BalanceAndOutputsService,
  ) { }

  ngOnInit() {
    // Get the synchronization status.
    this.subscriptionsGroup.push(this.blockchainService.progress.pipe(filter(response => !!response)).subscribe(response => {
      this.synchronizationInfoObtained = true;
      this.highestBlock = response.highestBlock;
      this.currentBlock = response.currentBlock;
      this.synchronizationPercentage = this.currentBlock && this.highestBlock ? (this.currentBlock / this.highestBlock) : 0;
      this.synchronized = response.synchronized;
    }));

    // Get the current price.
    this.subscriptionsGroup.push(this.priceService.price.subscribe(price => this.price = price));

    // Get the current balance.
    this.subscriptionsGroup.push(this.balanceAndOutputsService.walletsWithBalance.subscribe(wallets => {
      const addresses: AddressWithBalance[] = [];
      const alreadyAddedAddresses = new Map<string, boolean>();
      wallets.forEach(wallet => {
        wallet.addresses.forEach(address => {
          if (!alreadyAddedAddresses.has(address.address)) {
            addresses.push(address);
            alreadyAddedAddresses.set(address.address, true);
          }
        });
      });

      let coins = new BigNumber(0);
      let hours = new BigNumber(0);
      addresses.map(addr => {
        coins = coins.plus(addr.coins);
        hours = hours.plus(addr.hours);
      });
      this.coins = coins.toString();
      this.hours = hours.toString();

    }));

    // Know if there are pending transactions.
    this.subscriptionsGroup.push(this.balanceAndOutputsService.hasPendingTransactions.subscribe(hasPendingTxs => {
      this.hasPendingTxs = hasPendingTxs;
    }));

    // Know when the app gets the balance from the node.
    this.subscriptionsGroup.push(this.balanceAndOutputsService.firstFullUpdateMade.subscribe(firstFullUpdateMade => {
      this.balanceObtained = firstFullUpdateMade;
    }));
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }
}
