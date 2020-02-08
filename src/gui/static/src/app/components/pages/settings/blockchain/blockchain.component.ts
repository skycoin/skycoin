import { switchMap } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { SubscriptionLike, interval } from 'rxjs';
import { BlockchainService, BasicBlockInfo, CoinSupply } from '../../../../services/blockchain.service';
import { AppService } from '../../../../services/app.service';

@Component({
  selector: 'app-blockchain',
  templateUrl: './blockchain.component.html',
  styleUrls: ['./blockchain.component.scss'],
})
export class BlockchainComponent implements OnInit, OnDestroy {
  block: BasicBlockInfo;
  coinSupply: CoinSupply;

  private subscriptionsGroup: SubscriptionLike[] = [];

  constructor(
    public appService: AppService,
    private blockchainService: BlockchainService,
  ) { }

  ngOnInit() {
    this.subscriptionsGroup.push(interval(5000).pipe(switchMap(() => this.blockchainService.getLastBlock()))
      .subscribe(block => this.block = block));

    this.subscriptionsGroup.push(interval(5000).pipe(switchMap(() => this.blockchainService.getCoinSupply()))
      .subscribe(coinSupply => this.coinSupply = coinSupply),
    );
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }
}
