import { Component, Input, OnInit } from '@angular/core';
import { Transaction } from '../../../../../app.datatypes';
import { PriceService } from '../../../../../services/price.service';

@Component({
  selector: 'app-transaction-info',
  templateUrl: './transaction-info.component.html',
  styleUrls: ['./transaction-info.component.scss']
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
    this.transaction.hoursSent = this.transaction.outputs
      .filter(o => this.transaction.inputs.find(
        i => i[this.isPreview ? 'address' : 'owner'] !== o[this.isPreview ? 'address' : 'dst'])
      )
      .map(o => parseInt(o.hours, 10))
      .reduce((a, b) => a + b, 0);

    if (this.transaction.hoursSent === 0 && this.transaction.outputs.length === 1) {
      this.transaction.hoursSent = this.transaction.outputs[0].hours;
    }
  }

  toggleInputsOutputs(event) {
    event.preventDefault();

    this.showInputsOutputs = !this.showInputsOutputs;
  }
}
