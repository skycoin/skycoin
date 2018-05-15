import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { NormalTransaction } from '../../../../app.datatypes';

@Component({
  selector: 'app-transaction-detail',
  templateUrl: './transaction-detail.component.html',
  styleUrls: ['./transaction-detail.component.scss'],
})
export class TransactionDetailComponent {
  constructor(
    @Inject(MAT_DIALOG_DATA) public transaction: NormalTransaction,
    public dialogRef: MatDialogRef<TransactionDetailComponent>,
  ) {}

  closePopup() {
    this.dialogRef.close();
  }
}
