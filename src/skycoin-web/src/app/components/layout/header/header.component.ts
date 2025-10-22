import { Component, Input, OnDestroy, OnInit, NgZone } from '@angular/core';
import { Subscription } from 'rxjs';
import { BigNumber } from 'bignumber.js';
import { Observable } from 'rxjs';

import { PriceService } from '../../../services/price.service';
import { BalanceService, BalanceStates } from '../../../services/wallet/balance.service';
import { BlockchainService, ProgressEvent, ProgressStates } from '../../../services/blockchain.service';
import { ConnectionError } from '../../../enums/connection-error.enum';
import { CoinService } from '../../../services/coin.service';
import { BaseCoin } from '../../../coins/basecoin';
import { getTimeSinceLastBalanceUpdate } from '../../../utils';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Input() headline: string;

  coins: BigNumber = new BigNumber('0');
  hours: BigNumber;
  balance: string;
  hasPendingTxs: boolean;
  connectionError: ConnectionError = null;
  connectionErrorsList = ConnectionError;
  percentage: number;
  isBlockchainLoading = false;
  current: number;
  highest: number;
  currentCoin: BaseCoin;
  balanceObtained = false;
  timeSinceLastBalanceUpdate = 0;
  problemUpdatingBalance: boolean;
  synchronized = true;

  private price: number;
  private subscriptionsGroup: Subscription[] = [];
  private synchronizedSubscription: Subscription;

  get loading() {
    return this.isBlockchainLoading || !this.balanceObtained;
  }

  constructor(
    private priceService: PriceService,
    private balanceService: BalanceService,
    private blockchainService: BlockchainService,
    private coinService: CoinService,
    private _ngZone: NgZone
  ) {}

  ngOnInit() {
    this.subscriptionsGroup.push(this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => {
        this.currentCoin = coin;

        this.synchronized = true;
        if (this.synchronizedSubscription) {
          this.synchronizedSubscription.unsubscribe();
          this.synchronizedSubscription = null;
        }
      }));

    this.subscriptionsGroup.push(
      this.blockchainService.progress
        .filter(response => !!response)
        .subscribe(response => {
          this.updateBlockchainProgress(response);

          // Adding the code here prevents the warning from flashing if the wallet is synchronized. Also, adding the
          // subscription to this.subscription causes problems.
          if (response.currentBlock && !this.synchronizedSubscription) {
            this.synchronizedSubscription = this.blockchainService.synchronized.subscribe(value => this.synchronized = value);
          }
        })
    );

    this.subscriptionsGroup.push(
      this.priceService.price
        .subscribe(price => {
          this.price = price;
          this.calculateBalance();
        })
    );

    this.subscriptionsGroup.push(
      this.balanceService.totalBalance
        .subscribe(balance => {
          if (balance && balance.state === BalanceStates.Obtained) {
            this.coins = balance.balance.coins;
            this.hours = balance.balance.hours;
            this.balanceObtained = true;

            this.calculateBalance();
          }

          this.timeSinceLastBalanceUpdate = getTimeSinceLastBalanceUpdate(this.balanceService);
          this.problemUpdatingBalance = balance.state === BalanceStates.Error;
        })
    );

    this._ngZone.runOutsideAngular(() => {
      this.subscriptionsGroup.push(
        Observable.interval(5000).subscribe(() => this._ngZone.run(() => this.timeSinceLastBalanceUpdate = getTimeSinceLastBalanceUpdate(this.balanceService)))
      );
    });

    this.subscriptionsGroup.push(
      this.balanceService.hasPendingTransactions
        .subscribe(hasPendingTxs => this.hasPendingTxs = hasPendingTxs)
    );
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    if (this.synchronizedSubscription) {
      this.synchronizedSubscription.unsubscribe();
    }
  }

  private resetState() {
    this.coins = new BigNumber('0');
    this.price = null;
    this.balance = null;
    this.balanceObtained = false;
    this.isBlockchainLoading = false;
    this.percentage = null;
    this.current = null;
    this.highest = null;
  }

  private updateBlockchainProgress(response: ProgressEvent) {
    switch (response.state) {
      case ProgressStates.Restarting: {
        this.resetState();
        break;
      }
      case ProgressStates.Error: {
        this.setConnectionError(response.error);
        break;
      }
      case ProgressStates.Progress: {
        this.connectionError = null;
        this.isBlockchainLoading = response.highestBlock !== response.currentBlock;

        if (this.isBlockchainLoading) {
          this.highest = response.highestBlock;
          this.current = response.currentBlock;
        }

        this.percentage = response.currentBlock / response.highestBlock;
        break;
      }
    }
  }

  private calculateBalance() {
    if (this.price) {
      const balance = this.coins.multipliedBy(this.price).toNumber();
      this.balance = '$' + balance.toFixed(2) + ' ($' + (Math.round(this.price * 100) / 100).toFixed(2) + ')';
    }
  }

  private setConnectionError(error: ConnectionError) {
    if (!this.connectionError) {
      this.connectionError = error;
    }
  }
}
