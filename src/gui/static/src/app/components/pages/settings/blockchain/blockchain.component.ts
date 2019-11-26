import { Component, OnInit, OnDestroy } from '@angular/core';
import 'rxjs/add/operator/switchMap';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import { ISubscription } from 'rxjs/Subscription';
import { BlockchainService } from '../../../../services/blockchain.service';

@Component({
  selector: 'app-blockchain',
  templateUrl: './blockchain.component.html',
  styleUrls: ['./blockchain.component.scss'],
})
export class BlockchainComponent implements OnInit, OnDestroy {
  block: any;
  coinSupply: any;

  private subscriptionsGroup: ISubscription[] = [];

  constructor(
    private blockchainService: BlockchainService,
  ) { }

  ngOnInit() {
    this.subscriptionsGroup.push(IntervalObservable
      .create(5000)
      .switchMap(() => this.blockchainService.lastBlock())
      .subscribe(block => this.block = block));

    this.subscriptionsGroup.push(IntervalObservable
      .create(5000)
      .switchMap(() => this.blockchainService.coinSupply())
      .subscribe(coinSupply => this.coinSupply = coinSupply),
    );
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }
}
