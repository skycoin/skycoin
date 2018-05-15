import { Component, Input, OnInit } from '@angular/core';
import { PreviewTransaction, Transaction } from '../../../../../app.datatypes';
import { PriceService } from '../../../../../services/price.service';

@Component({
  selector: 'app-transaction-info',
  templateUrl: './transaction-info.component.html',
  styleUrls: ['./transaction-info.component.scss'],
})
export class TransactionInfoComponent implements OnInit {
  @Input() transaction: Transaction;
  @Input() isPreview: boolean;
  price: number;
  showInputsOutputs = false;

  constructor(private priceService: PriceService) {
    this.priceService.price.subscribe(price => this.price = price);
  }

  ngOnInit() {
    if (this.isPreview) {
      this.transaction.hoursSent = this.transaction.outputs
        .filter(o => o.address === (<PreviewTransaction> this.transaction).to)
        .map(o => parseInt(o.hours, 10))
        .reduce((a, b) => a + b, 0);
    }
  }

  toggleInputsOutputs(event) {
    event.preventDefault();

    this.showInputsOutputs = !this.showInputsOutputs;
  }
}
