import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { forkJoin, Subscription } from 'rxjs';

import { BlockchainService } from '../../../../services/blockchain.service';
import { CoinService } from '../../../../services/coin.service';
import { BaseCoin } from '../../../../coins/basecoin';

@Component({
    templateUrl: './blockchain.component.html',
    styleUrls: ['./blockchain.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class BlockchainComponent implements OnInit, OnDestroy {
  block: any;
  coinSupply: any;
  currentCoin!: BaseCoin;
  showError = false;

  private coinSubscription!: Subscription;
  private dataSubscription!: Subscription;

  constructor(
    private blockchainService: BlockchainService,
    private coinService: CoinService,
    private changeDetectorRef: ChangeDetectorRef,
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
            this.changeDetectorRef.markForCheck();
          },
          () => {
            this.showError = true;
            this.changeDetectorRef.markForCheck();
          });
        this.changeDetectorRef.markForCheck();
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
