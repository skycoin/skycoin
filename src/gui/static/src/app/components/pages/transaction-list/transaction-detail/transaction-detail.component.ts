import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { OldTransaction } from '../../../../services/wallet-operations/transaction-objects';

@Component({
  selector: 'app-transaction-detail',
  templateUrl: './transaction-detail.component.html',
  styleUrls: ['./transaction-detail.component.scss'],
})
export class TransactionDetailComponent {

  public static openDialog(dialog: MatDialog, transaction: OldTransaction): MatDialogRef<TransactionDetailComponent, any> {
    const config = new MatDialogConfig();
    config.data = transaction;
    config.autoFocus = false;
    config.width = '800px';

    return dialog.open(TransactionDetailComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public transaction: OldTransaction,
    public dialogRef: MatDialogRef<TransactionDetailComponent>,
  ) {}

  closePopup() {
    this.dialogRef.close();
  }
}
