import { mergeMap, delay } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { SubscriptionLike, of } from 'rxjs';

import { BlockchainService, BasicBlockInfo, CoinSupply } from '../../../../services/blockchain.service';
import { AppService } from '../../../../services/app.service';

/**
 * Shows the state of the the blockchain on the node.
 */
@Component({
  selector: 'app-blockchain',
  templateUrl: './blockchain.component.html',
  styleUrls: ['./blockchain.component.scss'],
})
export class BlockchainComponent implements OnInit, OnDestroy {
  block: BasicBlockInfo;
  coinSupply: CoinSupply;

  private operationSubscription: SubscriptionLike;

  // Time interval in which periodic data updates will be made.
  private readonly updatePeriod = 5 * 1000;
  // Time interval in which the periodic data updates will be restarted after an error.
  private readonly errorUpdatePeriod = 2 * 1000;

  constructor(
    public appService: AppService,
    private blockchainService: BlockchainService,
  ) { }

  ngOnInit() {
    this.startDataRefreshSubscription(0);
  }

  /**
   * Makes the page start updating the data periodically. If this function was called before,
   * the previous updating procedure is cancelled.
   * @param delayMs Delay before starting to update the data.
   */
  private startDataRefreshSubscription(delayMs: number) {
    this.removeOperationSubscription();

    this.operationSubscription = of(0).pipe(delay(delayMs), mergeMap(() => this.blockchainService.getLastBlock()), mergeMap(block => {
      this.block = block;

      return this.blockchainService.getCoinSupply();
    })).subscribe(coinSupply => {
      this.coinSupply = coinSupply;

      // Update again after a delay.
      this.startDataRefreshSubscription(this.updatePeriod);
    }, () => this.startDataRefreshSubscription(this.errorUpdatePeriod));
  }

  ngOnDestroy() {
    this.removeOperationSubscription();
  }

  private removeOperationSubscription() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
  }
}
