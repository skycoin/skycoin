import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Subscription } from 'rxjs';

import { PriceService } from '../../../../services/price.service';

@Component({
    selector: 'app-transaction-detail',
    templateUrl: './transaction-detail.component.html',
    styleUrls: ['./transaction-detail.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class TransactionDetailComponent implements OnInit, OnDestroy {

  price!: number;

  private priceSubscription!: Subscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public transaction: any,
    public dialogRef: MatDialogRef<TransactionDetailComponent>,
    private priceService: PriceService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {}

  ngOnInit() {
    this.priceSubscription = this.priceService.price.subscribe(price => {
      this.price = price;
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
  }

  closePopup() {
    this.dialogRef.close();
  }
}
