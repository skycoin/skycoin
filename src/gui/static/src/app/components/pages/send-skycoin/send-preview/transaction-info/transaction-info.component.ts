import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { PreviewTransaction, Transaction } from '../../../../../app.datatypes';
import { PriceService } from '../../../../../services/price.service';
import { ISubscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-transaction-info',
  templateUrl: './transaction-info.component.html',
  styleUrls: ['./transaction-info.component.scss'],
})
export class TransactionInfoComponent implements OnInit, OnDestroy {
  @Input() transaction: Transaction;
  @Input() isPreview: boolean;
  price: number;
  showInputsOutputs = false;

  private subscription: ISubscription;

  constructor(private priceService: PriceService) {
    this.subscription = this.priceService.price.subscribe(price => this.price = price);
  }

  ngOnInit() {
    if (this.isPreview) {
      this.transaction.hoursSent = this.transaction.outputs
        .filter(o => (<PreviewTransaction> this.transaction).to.find(addr => addr === o.address))
        .map(o => parseInt(o.hours, 10))
        .reduce((a, b) => a + b, 0);
    }
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  toggleInputsOutputs(event) {
    event.preventDefault();

    this.showInputsOutputs = !this.showInputsOutputs;
  }
}
