import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Transaction } from '../../../../app.datatypes';
import { PriceService } from '../../../../price.service';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-transaction-detail',
  templateUrl: './transaction-detail.component.html',
  styleUrls: ['./transaction-detail.component.scss']
})
export class TransactionDetailComponent implements OnInit, OnDestroy {

  price: number;

  private priceSubscription: Subscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public transaction: Transaction,
    public dialogRef: MatDialogRef<TransactionDetailComponent>,
    private priceService: PriceService,
  ) {}

  ngOnInit() {
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
  }

  closePopup() {
    this.dialogRef.close();
  }

  showOutput(output) {
    return !this.transaction.inputs.find(input => input.owner === output.dst);
  }
}
