import { Component, OnDestroy, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { forkJoin, Subscription } from 'rxjs';

import { BlockchainService } from '../../../../services/blockchain.service';
import { CoinService } from '../../../../services/coin.service';
import { BaseCoin } from '../../../../coins/basecoin';

@Component({
    templateUrl: './blockchain.component.html',
    styleUrls: ['./blockchain.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class BlockchainComponent implements OnInit, OnDestroy {
  block: any;
  coinSupply: any;
  currentCoin: BaseCoin;
  showError = false;

  private coinSubscription: Subscription;
  private dataSubscription: Subscription;

  constructor(
    private blockchainService: BlockchainService,
    private coinService: CoinService
  ) { }

  ngOnInit() {
    this.coinSubscription = this.coinService.currentCoin
      .subscribe((coin) => {
        this.currentCoin = coin;
        this.block = null;
        this.coinSupply = null;
        this.showError = false;

        this.closeDataSubscription();
        this.dataSubscription = forkJoin(
          this.blockchainService.lastBlock(),
          this.blockchainService.coinSupply())
          .subscribe(([block, coinSupply]) => {
            this.block = block;
            this.coinSupply = coinSupply;
          },
          () => this.showError = true);
      });
  }

  ngOnDestroy() {
    this.coinSubscription.unsubscribe();
    this.closeDataSubscription();
  }

  private closeDataSubscription() {
    if (this.dataSubscription && !this.dataSubscription.closed) {
      this.dataSubscription.unsubscribe();
    }
  }
}
