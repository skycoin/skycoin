import { switchMap } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { SubscriptionLike, interval } from 'rxjs';
import { BlockchainService } from '../../../../services/blockchain.service';
import { AppService } from '../../../../services/app.service';

@Component({
  selector: 'app-blockchain',
  templateUrl: './blockchain.component.html',
  styleUrls: ['./blockchain.component.scss'],
})
export class BlockchainComponent implements OnInit, OnDestroy {
  block: any;
  coinSupply: any;

  private subscriptionsGroup: SubscriptionLike[] = [];

  constructor(
    public appService: AppService,
    private blockchainService: BlockchainService,
  ) { }

  ngOnInit() {
    this.subscriptionsGroup.push(interval(5000).pipe(switchMap(() => this.blockchainService.lastBlock()))
      .subscribe(block => this.block = block));

    this.subscriptionsGroup.push(interval(5000).pipe(switchMap(() => this.blockchainService.coinSupply()))
      .subscribe(coinSupply => this.coinSupply = coinSupply),
    );
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }
}
