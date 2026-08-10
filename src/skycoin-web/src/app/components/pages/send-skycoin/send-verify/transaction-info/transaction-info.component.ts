import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';

import { Transaction } from '../../../../../app.datatypes';
import { PriceService } from '../../../../../services/price.service';
import { CoinService } from '../../../../../services/coin.service';
import { BaseCoin } from '../../../../../coins/basecoin';

@Component({
    selector: 'app-transaction-info',
    templateUrl: './transaction-info.component.html',
    styleUrls: ['./transaction-info.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class TransactionInfoComponent implements OnInit, OnDestroy {
  @Input() transaction!: Transaction;
  @Input() isPreview!: boolean;
  price!: number;
  showInputsOutputs = false;
  currentCoin!: BaseCoin;

  private subscription!: Subscription;

  constructor(
    private priceService: PriceService,
    private coinService: CoinService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {}

  get hoursText(): string {
    if (!this.transaction) {
      return '';
    }

    if (!this.isPreview) {
      if ((this.transaction as any).coinsMovedInternally) {
        return 'tx.hours-moved';
      } else if (this.transaction.balance!.isGreaterThan(0)) {
        return 'tx.hours-received';
      }
    }

    return 'tx.hours-sent';
  }

  ngOnInit() {
    this.subscription = this.priceService.price
      .subscribe(price => {
        this.price = price;
        this.changeDetectorRef.markForCheck();
      });

    this.currentCoin = this.coinService.currentCoin.value;
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  toggleInputsOutputs(event: any) {
    event.preventDefault();

    this.showInputsOutputs = !this.showInputsOutputs;
  }
}
