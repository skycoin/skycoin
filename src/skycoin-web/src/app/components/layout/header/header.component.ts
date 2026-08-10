import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, NgZone, OnDestroy, OnInit } from '@angular/core';
import { Subscription, interval, filter } from 'rxjs';
import { BigNumber } from 'bignumber.js';

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
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Input() headline!: string;

  coins: BigNumber = new BigNumber('0');
  hours!: BigNumber;
  balance!: string | null;
  hasPendingTxs!: boolean;
  connectionError: ConnectionError | null = null;
  connectionErrorsList = ConnectionError;
  percentage!: number | null;
  isBlockchainLoading = false;
  current!: number | null;
  highest!: number | null;
  currentCoin!: BaseCoin;
  balanceObtained = false;
  timeSinceLastBalanceUpdate = 0;
  problemUpdatingBalance!: boolean;
  synchronized = true;

  private price!: number | null;
  private subscriptionsGroup: Subscription[] = [];
  private synchronizedSubscription!: Subscription | null;

  get loading() {
    return this.isBlockchainLoading || !this.balanceObtained;
  }

  constructor(
    private priceService: PriceService,
    private balanceService: BalanceService,
    private blockchainService: BlockchainService,
    private coinService: CoinService,
    private _ngZone: NgZone,
    private changeDetectorRef: ChangeDetectorRef,
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
        this.changeDetectorRef.markForCheck();
      }));

    this.subscriptionsGroup.push(
      this.blockchainService.progress
        .pipe(filter(response => !!response))
        .subscribe(response => {
          this.updateBlockchainProgress(response);

          // Adding the code here prevents the warning from flashing if the wallet is synchronized. Also, adding the
          // subscription to this.subscription causes problems.
          if (response.currentBlock && !this.synchronizedSubscription) {
            this.synchronizedSubscription = this.blockchainService.synchronized.subscribe(value => {
              this.synchronized = value;
              this.changeDetectorRef.markForCheck();
            });
          }
          this.changeDetectorRef.markForCheck();
        })
    );

    this.subscriptionsGroup.push(
      this.priceService.price
        .subscribe(price => {
          this.price = price;
          this.calculateBalance();
          this.changeDetectorRef.markForCheck();
        })
    );

    this.subscriptionsGroup.push(
      this.balanceService.totalBalance
        .subscribe(balance => {
          if (balance && balance.state === BalanceStates.Obtained) {
            this.coins = balance.balance!.coins;
            this.hours = balance.balance!.hours;
            this.balanceObtained = true;

            this.calculateBalance();
          }

          this.timeSinceLastBalanceUpdate = getTimeSinceLastBalanceUpdate(this.balanceService);
          this.problemUpdatingBalance = balance.state === BalanceStates.Error;
          this.changeDetectorRef.markForCheck();
        })
    );

    this._ngZone.runOutsideAngular(() => {
      this.subscriptionsGroup.push(
        interval(5000).subscribe(() => {
          this._ngZone.run(() => this.timeSinceLastBalanceUpdate = getTimeSinceLastBalanceUpdate(this.balanceService));
          this.changeDetectorRef.markForCheck();
        })
      );
    });

    this.subscriptionsGroup.push(
      this.balanceService.hasPendingTransactions
        .subscribe(hasPendingTxs => {
          this.hasPendingTxs = hasPendingTxs;
          this.changeDetectorRef.markForCheck();
        })
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
        this.setConnectionError(response.error!);
        break;
      }
      case ProgressStates.Progress: {
        this.connectionError = null;
        this.isBlockchainLoading = response.highestBlock !== response.currentBlock;

        if (this.isBlockchainLoading) {
          this.highest = response.highestBlock!;
          this.current = response.currentBlock!;
        }

        this.percentage = response.currentBlock! / response.highestBlock!;
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
